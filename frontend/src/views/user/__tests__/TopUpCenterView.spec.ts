import { afterEach, describe, expect, it } from 'vitest'
import { mount, type VueWrapper } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { createMemoryHistory, createRouter } from 'vue-router'
import type { Component } from 'vue'
import zh from '@/i18n/locales/zh'

const pages = import.meta.glob<{ default: Component }>('../*View.vue')
const mounted: VueWrapper[] = []
async function renderCenter() {
  const load = pages['../TopUpCenterView.vue']
  expect(load, 'The card-shop recharge center must be available').toBeTypeOf('function')
  const component = (await load!()).default
  const router = createRouter({ history: createMemoryHistory(), routes: [
    { path: '/recharge', component },
    { path: '/redeem', component: { template: '<div>Redeem destination</div>' } },
  ] })
  await router.push('/recharge')
  await router.isReady()
  const wrapper = mount(component, { global: {
    plugins: [router, createI18n({ legacy: false, locale: 'zh', messages: { zh } })],
    stubs: { AppLayout: { template: '<main><slot /></main>' }, Icon: true },
  } })
  mounted.push(wrapper)
  return { wrapper, router }
}
afterEach(() => { mounted.splice(0).forEach(wrapper => wrapper.unmount()); document.documentElement.classList.remove('dark') })

describe('Card-shop recharge center', () => {
  it('opens the configured shop without exposing account identity or current URL secrets', async () => {
    window.history.replaceState({}, '', '/recharge?token=private-session&code=private-code')
    const { wrapper } = await renderCenter()
    const frame = wrapper.get('iframe')
    const url = new URL(frame.attributes('src'))
    expect(url.origin + url.pathname).toBe('https://wzyp.cn/shop/6B2BP3Q3')
    expect([...url.searchParams.keys()].sort()).toEqual(['lang', 'theme', 'ui_mode'])
    expect(url.searchParams.get('ui_mode')).toBe('embedded')
    expect(frame.attributes('referrerpolicy')).toBe('no-referrer')
    expect(wrapper.html()).not.toContain('private-session')
    expect(wrapper.html()).not.toContain('private-code')
    expect(frame.attributes('title')).toBeTruthy()
  })
  it('takes the buyer to the existing card redemption page', async () => {
    const { wrapper, router } = await renderCenter()
    await wrapper.get('a[href="/redeem"]').trigger('click')
    await new Promise(resolve => setTimeout(resolve, 0))
    expect(router.currentRoute.value.path).toBe('/redeem')
  })
  it('keeps a separate shop link available when embedding cannot load', async () => {
    const { wrapper } = await renderCenter()
    const link = wrapper.get('a[target="_blank"]')
    expect(link.attributes('href')).toBe('https://wzyp.cn/shop/6B2BP3Q3')
    expect(link.attributes('rel')).toContain('noopener')
    expect(link.attributes('rel')).toContain('noreferrer')
    await wrapper.get('iframe').trigger('error')
    expect(link.attributes('href')).toBe('https://wzyp.cn/shop/6B2BP3Q3')
    expect(link.isVisible()).toBe(true)
  })
  it('uses the current theme and updates the embedded shop when it changes', async () => {
    document.documentElement.classList.add('dark')
    const { wrapper } = await renderCenter()
    expect(new URL(wrapper.get('iframe').attributes('src')).searchParams.get('theme')).toBe('dark')
    document.documentElement.classList.remove('dark')
    await new Promise(resolve => setTimeout(resolve, 0))
    expect(new URL(wrapper.get('iframe').attributes('src')).searchParams.get('theme')).toBe('light')
  })
})
