<template>
  <AppLayout>
    <div class="top-up-center">
      <section class="card top-up-toolbar">
        <div class="min-w-0">
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('topUpCenter.buyCards') }}</h2>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('topUpCenter.instructions') }}</p>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <RouterLink to="/redeem" class="btn btn-primary">
            <Icon name="gift" size="sm" class="mr-1.5" />
            {{ t('topUpCenter.redeem') }}
          </RouterLink>
          <a :href="shopUrl" target="_blank" rel="noopener noreferrer" referrerpolicy="no-referrer" class="btn btn-secondary">
            <Icon name="externalLink" size="sm" class="mr-1.5" />
            {{ t('topUpCenter.openShop') }}
          </a>
        </div>
      </section>
      <section class="card top-up-shop" :aria-label="t('topUpCenter.shopTitle')">
        <iframe :src="embeddedUrl" :title="t('topUpCenter.shopTitle')" referrerpolicy="no-referrer" class="top-up-frame" />
      </section>
      <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('topUpCenter.loadHint') }}</p>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { detectTheme } from '@/utils/embedded-url'

const { t, locale } = useI18n()
const shopUrl = 'https://wzyp.cn/shop/6B2BP3Q3'
const theme = ref(detectTheme())
let observer: MutationObserver | undefined
// Card purchases do not need the user's identity, session token or source URL.
const embeddedUrl = computed(() => {
  const url = new URL(shopUrl)
  url.searchParams.set('ui_mode', 'embedded')
  url.searchParams.set('theme', theme.value)
  url.searchParams.set('lang', locale.value)
  return url.toString()
})
onMounted(() => {
  observer = new MutationObserver(() => { theme.value = detectTheme() })
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
})
onUnmounted(() => observer?.disconnect())
</script>

<style scoped>
.top-up-center { display: flex; min-width: 0; flex-direction: column; gap: 1rem; }
.top-up-toolbar { display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 1rem; padding: 1rem 1.25rem; }
.top-up-shop { overflow: hidden; }
.top-up-frame { display: block; width: 100%; height: max(640px, calc(100dvh - 280px)); border: 0; background: transparent; }
@media (max-width: 640px) {
  .top-up-toolbar { padding: 1rem; }
  .top-up-frame { height: max(680px, calc(100dvh - 300px)); }
}
</style>
