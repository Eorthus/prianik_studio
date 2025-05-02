import { defineNuxtPlugin } from 'nuxt/app'

export interface GsapMethods {
  initParallaxBanner: (
    banner: HTMLElement,
    image: HTMLElement
  ) => Promise<void>;
  initScrollReveal: (
    element: HTMLElement,
    options: Record<string, any>
  ) => Promise<void>;
  initStaggerAnimation: (
    elements: HTMLElement[],
    container: HTMLElement,
    options: Record<string, any>
  ) => Promise<void>;
  getGSAP: () => Promise<{ gsap: any; ScrollTrigger: any }>;
}

export default defineNuxtPlugin((nuxtApp) => {
  // Создаем переменную для хранения GSAP
  let gsap: typeof import('gsap').gsap | null = null;
  let ScrollTrigger: typeof import('gsap/ScrollTrigger').ScrollTrigger | null = null;
  
  // Создаем промис, который будет разрешен после загрузки GSAP
  const gsapLoaded: Promise<{ gsap: typeof import('gsap').gsap; ScrollTrigger: typeof import('gsap/ScrollTrigger').ScrollTrigger }> = new Promise(async (resolve) => {
    // Используем динамический импорт для ленивой загрузки GSAP
    // Это произойдет только тогда, когда функция будет вызвана
    const loadGSAP = async () => {
      if (!gsap) {
        const gsapModule = await import('gsap');
        gsap = gsapModule.gsap;
        
        // Загружаем плагин ScrollTrigger для анимаций на основе скролла
        const scrollTriggerModule = await import('gsap/ScrollTrigger');
        ScrollTrigger = scrollTriggerModule.ScrollTrigger;
        
        // Регистрируем плагин
        gsap.registerPlugin(ScrollTrigger);
        resolve({ gsap, ScrollTrigger });
      }
      return { gsap, ScrollTrigger };
    };
    
    await loadGSAP();
  });
  
  // Создаем функцию анимации параллакса для баннера
  const initParallaxBanner = async (bannerElement, imageElement) => {
    const { gsap, ScrollTrigger } = await gsapLoaded;
    
    // Создаем анимацию параллакса с усиленным эффектом
    gsap.to(imageElement, {
      yPercent: 50, // Увеличено с 30 до 50 для более заметного эффекта
      ease: "none",
      scrollTrigger: {
        trigger: bannerElement,
        start: "top top",
        end: "bottom top",
        scrub: 0.5, // Уменьшено для более быстрой реакции (было 1)
      }
    });
    
    // Добавим еще эффект затемнения при скролле для усиления эффекта
    const gradientOverlay = bannerElement.querySelector('.tw-bg-gradient-to-t');
    if (gradientOverlay) {
      gsap.to(gradientOverlay, {
        opacity: 0.9, // Начальная прозрачность 1, увеличивается до 0.9
        ease: "none",
        scrollTrigger: {
          trigger: bannerElement,
          start: "top top",
          end: "bottom top",
          scrub: true,
        }
      });
    }
  };
  
  // Создаем функцию для анимации появления элементов при скролле
  const initScrollReveal = async (elements, options = {}) => {
    const { gsap, ScrollTrigger } = await gsapLoaded;
    
    const defaults = {
      y: 50, // Смещение по вертикали
      opacity: 0, // Начальная прозрачность
      duration: 0.8, // Длительность анимации
      stagger: 0.1, // Задержка между анимациями для группы элементов
      ease: "power2.out", // Тип анимации
      scrollTrigger: {
        trigger: elements,
        start: "top 80%", // Начало анимации, когда 20% элемента видно
        toggleActions: "play none none none", // Поведение анимации при скролле
      }
    };
    
    // Объединяем пользовательские настройки с дефолтными
    const mergedOptions = { ...defaults, ...options };
    
    // Создаем анимацию
    return gsap.from(elements, mergedOptions);
  };
  
  // Создаем функцию для анимации стаггеров (элементов с задержкой)
  const initStaggerAnimation = async (elements, container, options = {}) => {
    const { gsap, ScrollTrigger } = await gsapLoaded;
    
    const defaults = {
      y: 30,
      opacity: 0,
      duration: 0.6,
      stagger: 0.1,
      ease: "power2.out",
      scrollTrigger: {
        trigger: container || elements,
        start: "top 80%",
      }
    };
    
    const mergedOptions = { ...defaults, ...options };
    
    return gsap.from(elements, mergedOptions);
  };
  
  // Функция для анимации сетки с эффектом каскада (для галереи)
  const initGridAnimation = async (gridItems, container, options = {}) => {
    const { gsap, ScrollTrigger } = await gsapLoaded;
    
    const defaults = {
      opacity: 0,
      y: 30,
      scale: 0.95,
      duration: 0.6,
      ease: "power2.out",
      stagger: {
        grid: [2, 3], // Автоматически определяется на основе CSS grid
        from: "start",
        amount: 0.6
      },
      scrollTrigger: {
        trigger: container,
        start: "top 85%",
      }
    };
    
    const mergedOptions = { ...defaults, ...options };
    //@ts-expect-error
    return gsap.from(gridItems, mergedOptions);
  };
  
  // Предоставляем набор функций через инъекцию зависимостей Nuxt
  nuxtApp.provide('gsap', {
    initParallaxBanner,
    initScrollReveal,
    initStaggerAnimation,
    initGridAnimation,
    getGSAP: async () => {
      return await gsapLoaded;
    }
  });
});