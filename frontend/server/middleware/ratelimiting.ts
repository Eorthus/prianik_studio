import { defineEventHandler, getRequestIP } from 'h3';

// Простая in-memory реализация (для production используйте Redis/другое хранилище)
const requestLimits = new Map();

export default defineEventHandler((event) => {
//   const ip = getRequestIP(event);
//   const path = event.node.req.url;
  
//   // Применяем rate limiting только к определенным путям (например, к API)
//   if (path?.startsWith('/api/')) {
//     const key = `${ip}-${path}`;
//     const now = Date.now();
//     const windowMs = 15 * 60 * 1000; // 15 минут
//     const maxRequests = 100; // максимум 100 запросов за период
    
//     const requestLog = requestLimits.get(key) || [];
    
//     // Очистка старых записей
//     const recentRequests = requestLog.filter(timestamp => now - timestamp < windowMs);
    
//     if (recentRequests.length >= maxRequests) {
//       event.node.res.statusCode = 429;
//       event.node.res.end('Too many requests, please try again later');
//       return;
//     }
    
//     // Добавляем текущий запрос
//     recentRequests.push(now);
//     requestLimits.set(key, recentRequests);
//   }
});