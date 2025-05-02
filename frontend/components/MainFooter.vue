<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<template>
  <footer class="tw-bg-white tw-border-t tw-border-gray-100 tw-py-8">
    <div class="tw-container tw-mx-auto tw-px-4">
      <div class="tw-flex tw-flex-col tw-items-center tw-text-center md:tw-text-left md:tw-flex-row tw-justify-between md:tw-items-start tw-gap-8">
        <!-- Юридическая информация -->
        <div class="tw-max-w-md">
          <h3 class="tw-text-lg tw-font-medium tw-mb-4 tw-text-gray-800">{{ $t('contacts.company_name') }}</h3>
          <p class="tw-text-sm tw-text-gray-600 tw-mb-2">{{ $t('footer.legal_entity') }}</p>
          <p class="tw-text-sm tw-text-gray-600 tw-mb-2">{{ $t('footer.legal_id') }}: 123456789012345</p>
          <p class="tw-text-sm tw-text-gray-600">{{ $t('footer.company') }} {{ new Date().getFullYear() }} {{ $t('footer.all_rights') }}</p>
        </div>

        <!-- Контактная информация -->
        <div>
          <h3 class="tw-text-lg tw-font-medium tw-mb-4 tw-text-gray-800">{{ $t('contacts.title') }}</h3>
          <div class="tw-flex tw-flex-col tw-gap-2">
            <a href="mailto:info@creality-workshop.ru" class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-500">
              {{ $t('contacts.company_email_prefix') }}@{{ $t('contacts.company_email_domain') }}
            </a>
            <a href="tel:+79001234567" class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-500">
              {{ $t('contacts.company_phone') }}
            </a>
            <p class="tw-text-sm tw-text-gray-600">{{ $t('contacts.work_hours') }}: {{ $t('footer.working_hours_short') }}</p>
          </div>
        </div>

        <!-- Социальные сети -->
        <div>
          <h3 class="tw-text-lg tw-font-medium tw-mb-4 tw-text-gray-800">{{ $t('contacts.social_media') }}</h3>
          <div class="tw-flex tw-gap-4">
            <a href="#" class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-500" aria-label="VK">
              <svg class="tw-w-6 tw-h-6" viewBox="0 0 24 24" fill="currentColor">
                <path d="M21.579 6.855c.14-.465 0-.806-.662-.806h-2.193c-.558 0-.813.295-.953.619 0 0-1.115 2.719-2.695 4.482-.51.513-.741.675-1.021.675-.139 0-.341-.162-.341-.627V6.855c0-.558-.161-.806-.626-.806H9.642c-.348 0-.558.258-.558.504 0 .528.788.65.869 2.138v3.228c0 .707-.128.836-.406.836-.742 0-2.546-2.725-3.616-5.846-.209-.606-.42-.854-.98-.854H2.759c-.626 0-.752.295-.752.619 0 .582.742 3.462 3.461 7.271 1.812 2.601 4.363 4.011 6.687 4.011 1.393 0 1.565-.313 1.565-.853v-1.966c0-.626.133-.752.574-.752.324 0 .882.164 2.183 1.417 1.486 1.486 1.732 2.153 2.567 2.153h2.192c.626 0 .939-.313.759-.931-.197-.615-.907-1.51-1.849-2.569-.512-.604-1.277-1.254-1.51-1.579-.325-.419-.231-.604 0-.976.001.001 2.672-3.761 2.95-5.04z" />
              </svg>
            </a>
            <a href="#" class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-500" aria-label="Telegram">
              <svg class="tw-w-6 tw-h-6" viewBox="0 0 24 24" fill="currentColor">
                <path d="M22.05 1.577c-.393-.016-.784.08-1.117.235-.484.186-4.92 1.902-9.41 3.64-2.26.873-4.518 1.746-6.256 2.415-1.737.67-3.045 1.168-3.114 1.192-.46.16-1.082.362-1.61.984-.133.155-.267.354-.335.628s-.038.622.095.895c.265.547.714.773 1.244.976.53.204 1.525.573 2.021.76.496.187 1.091.413 1.776.607.715.19 1.693-.367 1.976-1.09.284-.722.552-2.014.808-2.64.256-.625.6-.463.6-.463l6.516 2.09s.798.242 1.134-.202c.334-.444.105-.958.105-.958l-4.003-6.116s-.379-.585-1.004-.724c-.315-.07-.541-.116-.921-.13zm-3.423 8.93c-.088.16-.177.312-.208.425-.243.13-3.52 2.818-4.92 3.138-.33.477-.713.86-1.25.674-.538-.186-2.433-1.02-3.45-1.411-1.8-.693-3.43-1.32-3.43-1.32s-.566-.215-.398-.798c1.398-.913 4.247-1.75 5.892-2.344 1.645-.594 3.2-1.131 3.2-1.131s.636-.259 1.005.222c.368.48 3.43 2.824 3.43 2.824s.256.28.17.44z" />
              </svg>
            </a>
            <a href="#" class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-500" aria-label="Instagram">
              <svg class="tw-w-6 tw-h-6" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z" />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>