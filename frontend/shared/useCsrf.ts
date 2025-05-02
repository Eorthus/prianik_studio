import { ref } from 'vue';

export function useCsrf() {
  const csrfToken = ref('');
  
  // Генерирует рандомный токен
  const generateToken = () => {
    const random = Math.random().toString(36).substring(2, 15) + 
                  Math.random().toString(36).substring(2, 15);
    csrfToken.value = random;
    
    // Сохраняем в localStorage (в реальном приложении лучше использовать HttpOnly cookie)
    localStorage.setItem('csrf_token', csrfToken.value);
    
    return csrfToken.value;
  };
  
  // Проверяет валидность токена
  const validateToken = (token) => {
    const storedToken = localStorage.getItem('csrf_token');
    return token === storedToken;
  };
  
  return {
    csrfToken,
    generateToken,
    validateToken
  };
}