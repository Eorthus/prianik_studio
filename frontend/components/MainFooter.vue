<script setup>
import { useI18n } from "vue-i18n";
import WhatsAppIcon from "./icons/WhatsAppIcon.vue";
import InstagramIcon from "./icons/InstagramIcon.vue";
const { t } = useI18n();
</script>

<template>
  <footer class="tw-bg-white tw-border-t tw-border-gray-100 tw-py-8">
    <div class="tw-container tw-mx-auto tw-px-4">
      <div
        class="tw-flex tw-flex-col tw-items-center tw-text-center md:tw-text-left md:tw-flex-row tw-justify-between md:tw-items-start tw-gap-8"
      >
        <!-- Юридическая информация -->
        <div class="tw-max-w-md">
          <h3 class="tw-text-lg tw-font-medium tw-mb-4 tw-text-gray-800">
            {{ $t("contacts.company_name") }}
          </h3>
          <!-- <p class="tw-text-sm tw-text-gray-600 tw-mb-2">{{ $t('footer.legal_entity') }}</p>
          <p class="tw-text-sm tw-text-gray-600 tw-mb-2">{{ $t('footer.legal_id') }}: 123456789012345</p> -->
          <p class="tw-text-sm tw-text-gray-600">
            {{ $t("footer.company") }} {{ new Date().getFullYear() }}
            {{ $t("footer.all_rights") }}
          </p>
        </div>

        <!-- Контактная информация -->
        <div>
          <h3 class="tw-text-lg tw-font-medium tw-mb-4 tw-text-gray-800">
            {{ $t("contacts.title") }}
          </h3>
          <div class="tw-flex tw-flex-col tw-gap-2">
            <a
              href="mailto:info@creality-workshop.ru"
              class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-500"
            >
              {{ $t("contacts.company_email_prefix") }}@{{
                $t("contacts.company_email_domain")
              }}
            </a>
            <a
              href="tel:+79001234567"
              class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-500"
            >
              {{ $t("contacts.company_phone") }}
            </a>
            <p class="tw-text-sm tw-text-gray-600">
              {{ $t("contacts.work_hours") }}:
              {{ $t("footer.working_hours_short") }}
            </p>
          </div>
        </div>

        <!-- Социальные сети -->
        <div>
          <h3 class="tw-text-lg tw-font-medium tw-mb-4 tw-text-gray-800">
            {{ $t("contacts.social_media") }}
          </h3>
          <div class="tw-flex tw-gap-4">
            <a
              :href="`https://wa.me/${$t('contacts.company_phone').replace(
                /[^\d]/g,
                ''
              )}`"
              target="_blank"
              rel="noopener noreferrer"
              class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-500"
              aria-label="WhatsApp"
            >
              <WhatsAppIcon class="tw-w-6 tw-h-6" />
            </a>
            <a
              href="https://instagram.com/prianik_studio"
              target="_blank"
              rel="noopener noreferrer"
              class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-500"
              aria-label="Instagram"
            >
              <InstagramIcon class="tw-w-6 tw-h-6" />
            </a>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>
