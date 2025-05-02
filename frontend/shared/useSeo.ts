import { useRuntimeConfig } from "nuxt/app";
import type { RouteLocationNormalizedLoaded } from "vue-router";

export interface SeoMetaOptions {
  title?: string;
  description?: string;
  image?: string;
  url?: string;
  type?: string;
  twitterCard?: string;
}

// Чистая функция для получения заголовка
export function getPageTitle(
  t: any,
  routeName: string,
  customTitle?: string
): string {
  if (customTitle) return customTitle;

  // Удаляем суффикс языка из имени маршрута (например, gallery__es -> gallery)
  const baseRouteName = routeName.split("__")[0];

  // Правильно используем t() для перевода
  // Избегаем динамических ключей, используем явные условия для каждой страницы
  let translatedTitle = "";

  switch (baseRouteName) {
    case "index":
      translatedTitle = t("meta.home");
      break;
    case "catalog":
    case "catalog-id":
      translatedTitle = t("meta.catalog");
      break;
    case "gallery":
      translatedTitle = t("meta.gallery");
      break;
    case "contacts":
      translatedTitle = t("meta.contacts");
      break;
    case "cart":
      translatedTitle = t("meta.cart");
      break;
    default:
      translatedTitle = t("meta.home");
  }

  // Возвращаем перевод, а не ключ
  return translatedTitle || "Prianik Studio";
}

// Чистая функция для получения описания
export function getPageDescription(
  t: any,
  routeName: string,
  customDescription?: string
): string {
  if (customDescription) return customDescription;

  // Удаляем суффикс языка из имени маршрута
  const baseRouteName = routeName.split("__")[0];

  // Правильно используем t() для перевода
  // Избегаем динамических ключей, используем явные условия для каждой страницы
  let translatedDescription = "";

  switch (baseRouteName) {
    case "index":
      translatedDescription = t("meta.description.home");
      break;
    case "catalog":
    case "catalog-id":
      translatedDescription = t("meta.description.catalog");
      break;
    case "gallery":
      translatedDescription = t("meta.description.gallery");
      break;
    case "contacts":
      translatedDescription = t("meta.description.contacts");
      break;
    case "cart":
      translatedDescription = t("meta.description.cart");
      break;
    default:
      translatedDescription = t("meta.description.home");
  }

  // Возвращаем перевод, а не ключ
  return (
    translatedDescription || "Лазерная гравировка и 3D печать любой сложности"
  );
}

// Чистая функция для получения URL изображения
export function getPageImage(locale: string, customImage?: string): string {
  if (customImage) return customImage;
  return `../assets/img/prianik_og.png`;
  // Правильный путь к изображению в зависимости от языка
  // switch(locale) {
  //   case 'es':
  //     return `../assets/img/prianik_og.png`;
  //   case 'en':
  //     return `../assets/img/prianik_og.png`;
  //   case 'ru':
  //     return `../assets/img/prianik_og.png`;
  //   default:
  //     return `../assets/img/prianik_og.png`;
  // }
}

// Чистая функция для получения URL страницы
export function getPageUrl(
  route: RouteLocationNormalizedLoaded,
  customUrl?: string
): string {
  const API_BASE_URL = useRuntimeConfig().public.apiBaseUrl;
  return customUrl || `${API_BASE_URL}${route.fullPath}`;
}

// Функция для создания альтернативных языковых ссылок
export function getAlternateLinks(
  route: RouteLocationNormalizedLoaded,
  currentLocale: string
): Array<{ rel: string; hrefLang: string; href: string }> {
  const API_BASE_URL = useRuntimeConfig().public.apiBaseUrl;
  const locales = ["es", "en", "ru"];
  const defaultLocale = "es";

  return locales.map((locale) => {
    // Формируем путь для каждого языка
    let path = route.fullPath;

    if (locale !== currentLocale) {
      // Если текущий маршрут не имеет языкового префикса (он на языке по умолчанию)
      const hasLocalePrefix = locales.some((loc) =>
        route.fullPath.startsWith(`/${loc}`)
      );

      if (!hasLocalePrefix && currentLocale === defaultLocale) {
        // Мы находимся на маршруте по умолчанию без префикса
        if (locale !== defaultLocale) {
          // Добавляем префикс для неосновных языков
          path = `/${locale}${path}`;
        }
      } else if (locale === defaultLocale) {
        // Для языка по умолчанию убираем любой текущий языковой префикс
        path = path.replace(new RegExp(`^/${currentLocale}`), "");
      } else if (currentLocale === defaultLocale) {
        // Текущий язык - язык по умолчанию (без префикса), добавляем префикс для других языков
        path = `/${locale}${path}`;
      } else {
        // Заменяем префикс текущего языка на префикс целевого языка
        path = path.replace(new RegExp(`^/${currentLocale}`), `/${locale}`);
      }
    }

    return {
      rel: "alternate",
      hrefLang: locale,
      href: `${API_BASE_URL}${path}`,
    };
  });
}

// Функция для формирования объекта метаданных
export function getSeoMetadata(
  route: RouteLocationNormalizedLoaded,
  t: any,
  locale: string,
  options: SeoMetaOptions = {}
) {
  const routeName = route.name?.toString() || "index";

  // Получаем данные с помощью чистых функций
  const title = getPageTitle(t, routeName, options.title);
  const description = getPageDescription(t, routeName, options.description);
  const image = getPageImage(locale, options.image);
  const pageUrl = getPageUrl(route, options.url);
  const type = options.type || "website";
  const twitterCard = options.twitterCard || "summary_large_image";
  const alternateLinks = getAlternateLinks(route, locale);

  // Возвращаем метаданные
  return {
    title,
    description,
    image,
    pageUrl,
    locale,
    type,
    twitterCard,
    alternateLinks,
  };
}
