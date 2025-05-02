/**
 * Экранирует HTML-специальные символы для предотвращения XSS
 */
export function escapeHtml(unsafe: string): string {
    return unsafe
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#039;");
  }
  
  /**
   * Проверяет ввод на вредоносные паттерны
   */
  export function containsMaliciousPatterns(input: string): boolean {
    const maliciousPatterns = [
      /<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi,
      /javascript:/gi,
      /onerror=/gi,
      /onload=/gi,
      /onclick=/gi,
      /eval\(/gi
    ];
    
    return maliciousPatterns.some(pattern => pattern.test(input));
  }
  
  /**
   * Санитизация пользовательского ввода
   */
  export function sanitizeInput(input: string): string {
    if (typeof input !== 'string') return '';
    
    // Обрезаем пробелы
    let sanitized = input.trim();
    
    // Экранируем HTML
    sanitized = escapeHtml(sanitized);
    
    return sanitized;
  }