import { defineEventHandler } from 'h3';

export default defineEventHandler(async (event) => {
  // Устанавливаем заголовок Content-Type
  event.node.res.setHeader('Content-Type', 'application/xml');
  
  // Базовый URL сайта
  const siteUrl = 'https://creality-workshop.com';
  
  // Языки сайта
  const locales = ['es', 'en', 'ru'];
  const defaultLocale = 'es';
  
  // Статические маршруты
  const staticRoutes = [
    '/',
    '/catalog',
    '/gallery',
    '/contacts',
    '/cart'
  ];
  
  // Динамические маршруты (для примера - ID товаров)
  // В реальном приложении здесь могла бы быть логика получения 
  // всех ID товаров из базы данных или API
  const productIds = Array.from({ length: 15 }, (_, i) => i + 1);
  const productRoutes = productIds.map(id => `/catalog/${id}`);
  
  // Объединяем все маршруты
  const allRoutes = [...staticRoutes, ...productRoutes];
  
  // Генерируем XML
  let sitemap = '<?xml version="1.0" encoding="UTF-8"?>\n';
  sitemap += '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">\n';
  
  // Добавляем все маршруты для всех локалей
  for (const route of allRoutes) {
    for (const locale of locales) {
      const path = locale === defaultLocale ? route : `/${locale}${route}`;
      const url = `${siteUrl}${path}`;
      
      sitemap += '  <url>\n';
      sitemap += `    <loc>${url}</loc>\n`;
      
      // Добавляем альтернативные языковые версии
      for (const altLocale of locales) {
        const altPath = altLocale === defaultLocale ? route : `/${altLocale}${route}`;
        const altUrl = `${siteUrl}${altPath}`;
        
        sitemap += `    <xhtml:link rel="alternate" hreflang="${altLocale}" href="${altUrl}" />\n`;
      }
      
      // Приоритет и частота обновления
      const priority = route === '/' ? '1.0' : (route.includes('/catalog/') ? '0.8' : '0.7');
      sitemap += `    <priority>${priority}</priority>\n`;
      sitemap += '    <changefreq>weekly</changefreq>\n';
      sitemap += `    <lastmod>${new Date().toISOString()}</lastmod>\n`;
      sitemap += '  </url>\n';
    }
  }
  
  sitemap += '</urlset>';
  
  return sitemap;
});