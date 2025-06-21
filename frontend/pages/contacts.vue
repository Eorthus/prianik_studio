<script setup lang="ts">
import { useI18n } from "vue-i18n";
import FAQAccordion from "~/components/contacts/FAQAccordion.vue";
import DestinationIcon from "~/components/icons/DestinationIcon.vue";
import PhoneIcon from "~/components/icons/PhoneIcon.vue";
import EmailIcon from "~/components/icons/EmailIcon.vue";
import ScheduleIcon from "~/components/icons/ScheduleIcon.vue";
import ContactFormHandler from "~/components/contacts/ContactFormHandler.vue";
import WhatsAppIcon from "~/components/icons/WhatsAppIcon.vue";
import InstagramIcon from "~/components/icons/InstagramIcon.vue";
const { t } = useI18n();

// Полностью локализованные данные о компании
const companyInfo = {
  name: t("contacts.company_name"),
  address: t("contacts.company_address"),
  phone: t("contacts.company_phone"),
  email: ` ${t("contacts.company_email_prefix")}@${t(
    "contacts.company_email_domain"
  )}`,
  workHours: t("contacts.company_hours"),
};

// Данные для FAQ - получаем из локализации
const faqItems = [
  {
    question: t("contacts.faq.custom_order.question"),
    answer: t("contacts.faq.custom_order.answer"),
  },
  {
    question: t("contacts.faq.production_time.question"),
    answer: t("contacts.faq.production_time.answer"),
  },
  {
    question: t("contacts.faq.delivery.question"),
    answer: t("contacts.faq.delivery.answer"),
  },
  {
    question: t("contacts.faq.materials.question"),
    answer: t("contacts.faq.materials.answer"),
  },
];
</script>

<template>
  <div class="tw-bg-gray-50 tw-pt-16">
    <h1
      class="tw-text-3xl tw-font-bold tw-text-gray-800 tw-mb-8 tw-text-center"
    >
      {{ $t("contacts.title") }}
    </h1>
    <!-- Информация о компании -->
    <div
      class="tw-grid tw-grid-cols-1 md:tw-grid-cols-2 tw-gap-12 tw-p-8 tw-pt-0 tw-container tw-mx-auto"
    >
      <div class="tw-p-8">
        <h2 class="tw-text-xl tw-font-semibold tw-mb-6">
          {{ $t("contacts.our_contacts") }}
        </h2>

        <div class="tw-space-y-4">
          <div class="tw-flex tw-items-start">
            <div class="tw-flex-shrink-0 tw-text-gray-600">
              <DestinationIcon class="tw-w-6 tw-h-6" />
            </div>
            <div class="tw-ml-4">
              <p class="tw-text-sm tw-font-medium tw-text-gray-800">
                {{ companyInfo.name }}
              </p>
              <p class="tw-text-sm tw-text-gray-600">
                {{ companyInfo.address }}
              </p>
            </div>
          </div>

          <div class="tw-flex tw-items-start">
            <div class="tw-flex-shrink-0 tw-text-gray-600">
              <PhoneIcon class="tw-w-6 tw-h-6" />
            </div>
            <div class="tw-ml-4">
              <p class="tw-text-sm tw-font-medium tw-text-gray-800">
                {{ $t("contacts.phone") }}
              </p>
              <a
                href="tel:+79001234567"
                class="tw-text-sm tw-text-gray-600 tw-transition-colors hover:tw-text-gray-800"
              >
                {{ companyInfo.phone }}
              </a>
            </div>
          </div>

          <div class="tw-flex tw-items-start">
            <div class="tw-flex-shrink-0 tw-text-gray-600">
              <EmailIcon class="tw-w-6 tw-h-6" />
            </div>
            <div class="tw-ml-4">
              <p class="tw-text-sm tw-font-medium tw-text-gray-800">
                {{ $t("contacts.email") }}
              </p>
              <a
                :href="'mailto:' + companyInfo.email"
                class="tw-text-sm tw-text-gray-600 tw-transition-colors hover:tw-text-gray-800"
              >
                {{ companyInfo.email }}
              </a>
            </div>
          </div>

          <div class="tw-flex tw-items-start">
            <div class="tw-flex-shrink-0 tw-text-gray-600">
              <ScheduleIcon class="tw-w-6 tw-h-6" />
            </div>
            <div class="tw-ml-4">
              <p class="tw-text-sm tw-font-medium tw-text-gray-800">
                {{ $t("contacts.work_hours") }}
              </p>
              <p class="tw-text-sm tw-text-gray-600">
                {{ companyInfo.workHours }}
              </p>
            </div>
          </div>
        </div>
      </div>
      <!-- Социальные сети -->
      <div class="tw-p-8">
        <h3 class="tw-text-xl tw-font-semibold tw-text-gray-800 tw-mb-4">
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
            class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-800"
            aria-label="WhatsApp"
          >
            <WhatsAppIcon class="tw-w-6 tw-h-6" />
          </a>
          <a
            href="https://instagram.com/prianik_studio"
            target="_blank"
            rel="noopener noreferrer"
            class="tw-text-gray-600 tw-transition-colors tw-duration-300 hover:tw-text-gray-800"
            aria-label="Instagram"
          >
            <InstagramIcon class="tw-w-6 tw-h-6" />
          </a>
        </div>
      </div>
    </div>
  </div>

  <!-- Форма обратной связи -->
  <section class="tw-py-16 tw-bg-white">
    <ContactFormHandler />
  </section>

  <!-- FAQ - Часто задаваемые вопросы -->
  <div class="tw-py-16 tw-bg-gray-50">
    <h2
      class="tw-text-xl tw-font-semibold tw-text-gray-800 tw-mb-6 tw-text-center"
    >
      {{ $t("contacts.faq_title") }}
    </h2>

    <!-- Используем компонент FaqAccordion -->
    <FAQAccordion :faqs="faqItems" />
  </div>
</template>

<style scoped>
/* Стиль для карты с фиксированным соотношением сторон */
.tw-aspect-w-16 {
  position: relative;
  padding-bottom: 56.25%; /* 16:9 */
}

.tw-aspect-w-16 > * {
  position: absolute;
  height: 100%;
  width: 100%;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
}
</style>
