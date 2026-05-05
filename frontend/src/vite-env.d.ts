/// <reference types="vite/client" />

interface Window {
  go?: {
    main?: {
      App?: Record<string, (...args: unknown[]) => Promise<unknown>>;
    };
  };
  runtime?: {
    Quit?: () => void;
    WindowHide?: () => void;
    WindowShow?: () => void;
  };
}
