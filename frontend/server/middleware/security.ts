export default defineEventHandler((event) => {
    const headers = event.node.res.getHeaders();
    
    // Защита от XSS
    event.node.res.setHeader('X-XSS-Protection', '1; mode=block');
    
    // Предотвращение кликджекинга
    event.node.res.setHeader('X-Frame-Options', 'SAMEORIGIN');
    
    // Запрет браузеру угадывать тип контента
    event.node.res.setHeader('X-Content-Type-Options', 'nosniff');
    
    // Базовая CSP для защиты от XSS
    event.node.res.setHeader(
      'Content-Security-Policy',
      "default-src 'self'; script-src 'self' https://cdnjs.cloudflare.com; style-src 'self' https://fonts.googleapis.com 'unsafe-inline'; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https://s3.stroi-news.ru; connect-src 'self'"
    );
    
    // Разрешаем загрузку только по HTTPS (если сайт работает по HTTPS)
    event.node.res.setHeader(
      'Strict-Transport-Security',
      'max-age=31536000; includeSubDomains'
    );
    
    // Параметры безопасности для куки
    event.node.res.setHeader(
      'Set-Cookie',
      'Path=/; HttpOnly; Secure; SameSite=Strict'
    );
  });