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
export function getPageTitle(t: any, routeName: string, customTitle?: string): string {
  if (customTitle) return customTitle;

    // Удаляем суффикс языка из имени маршрута (например, gallery__es -> gallery)
    const baseRouteName = routeName.split('__')[0];
  
  // Правильно используем t() для перевода
  // Избегаем динамических ключей, используем явные условия для каждой страницы
  let translatedTitle = '';
  
  switch(baseRouteName) {
    case 'index':
      translatedTitle = t('meta.home');
      break;
    case 'catalog':
    case 'catalog-id':
      translatedTitle = t('meta.catalog');
      break;
    case 'gallery':
      translatedTitle = t('meta.gallery');
      break;
    case 'contacts':
      translatedTitle = t('meta.contacts');
      break;
    case 'cart':
      translatedTitle = t('meta.cart');
      break;
    default:
      translatedTitle = t('meta.home');
  }
  
  // Возвращаем перевод, а не ключ
  return translatedTitle || 'Prianik Studio';
}

// Чистая функция для получения описания
export function getPageDescription(t: any, routeName: string, customDescription?: string): string {
  if (customDescription) return customDescription;

   // Удаляем суффикс языка из имени маршрута
   const baseRouteName = routeName.split('__')[0];
  
  // Правильно используем t() для перевода
  // Избегаем динамических ключей, используем явные условия для каждой страницы
  let translatedDescription = '';
  
  switch(baseRouteName) {
    case 'index':
      translatedDescription = t('meta.description.home');
      break;
    case 'catalog':
    case 'catalog-id':
      translatedDescription = t('meta.description.catalog');
      break;
    case 'gallery':
      translatedDescription = t('meta.description.gallery');
      break;
    case 'contacts':
      translatedDescription = t('meta.description.contacts');
      break;
    case 'cart':
      translatedDescription = t('meta.description.cart');
      break;
    default:
      translatedDescription = t('meta.description.home');
  }
  
  // Возвращаем перевод, а не ключ
  return translatedDescription || 'Лазерная гравировка и 3D печать любой сложности';
}

// Чистая функция для получения URL изображения
export function getPageImage(locale: string, customImage?: string): string {
  if (customImage) return customImage;
  
  const baseUrl = "https://creality-workshop.com";
  
  // Правильный путь к изображению в зависимости от языка
  switch(locale) {
    case 'es':
      return `${baseUrl}/images/og-image-es.jpg`;
    case 'en':
      return `${baseUrl}/images/og-image-en.jpg`;
    case 'ru':
      return `${baseUrl}/images/og-image-ru.jpg`;
    default:
      return `${baseUrl}/images/og-image-default.jpg`;
  }
}

// Чистая функция для получения URL страницы
export function getPageUrl(route: RouteLocationNormalizedLoaded, customUrl?: string): string {
  const baseUrl = "https://creality-workshop.com";
  return customUrl || `${baseUrl}${route.fullPath}`;
}

// Функция для создания альтернативных языковых ссылок
export function getAlternateLinks(route: RouteLocationNormalizedLoaded, currentLocale: string): Array<{ rel: string, hrefLang: string, href: string }> {
  const baseUrl = "https://creality-workshop.com";
  const locales = ['es', 'en', 'ru'];
  const defaultLocale = 'es';
  
  return locales.map(locale => {
    // Формируем путь для каждой локали
    let path = route.fullPath;
    
    if (locale !== currentLocale) {
      if (locale === defaultLocale) {
        // Для дефолтной локали убираем префикс текущей локали
        path = path.replace(new RegExp(`^/${currentLocale}`), '');
      } else if (currentLocale === defaultLocale) {
        // Текущая локаль дефолтная, добавляем префикс для других локалей
        path = `/${locale}${path}`;
      } else {
        // Заменяем префикс текущей локали на префикс целевой локали
        path = path.replace(new RegExp(`^/${currentLocale}`), `/${locale}`);
      }
    }
    
    return {
      rel: "alternate",
      hrefLang: locale,
      href: `${baseUrl}${path}`
    };
  });
}

// Функция для формирования объекта метаданных
export function getSeoMetadata(route: RouteLocationNormalizedLoaded, t: any, locale: string, options: SeoMetaOptions = {}) {
  const routeName = route.name?.toString() || 'index';
  
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
    alternateLinks
  };
}