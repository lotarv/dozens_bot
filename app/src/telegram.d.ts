interface TelegramWebApp {
    initData: string;
    initDataUnsafe: any;
    MainButton: {
      setText(text: string): void;
      show(): void;
      hide(): void;
      onClick(callback: () => void): void;
    };
    disableVerticalSwipes(): void;
    close(): void;
    // Добавьте другие методы и свойства Telegram Web App API, которые используете
  }
  
  interface Window {
    Telegram: {
      WebApp: TelegramWebApp;
    };
  }