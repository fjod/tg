import React, { useEffect, useState } from 'react';
import Header from './components/Header.jsx';
import TagList from './components/TagList.jsx';
import telegramApp from './utils/telegram.js';
import './styles.css';

function App() {
  const [isInitialized, setIsInitialized] = useState(false);
  const [error, setError] = useState(null);
  const [theme, setTheme] = useState({});

  useEffect(() => {
    initializeApp();
  }, []);

  const initializeApp = async () => {
    try {
      console.log('Initializing Telegram Mini-App...');
      
      // Initialize Telegram WebApp
      const initialized = telegramApp.init();
      
      if (!initialized) {
        console.warn('Running outside Telegram - using mock data for development');
        setError('This app should be opened from Telegram');
      }

      // Check if we have user authentication
      const user = telegramApp.getUser();
      const authData = telegramApp.getAuthData();

      console.log('App initialization:', {
        initialized,
        hasUser: !!user,
        hasAuthData: !!authData,
        isInTelegram: telegramApp.isInTelegram(),
      });

      if (!telegramApp.isInTelegram() && process.env.NODE_ENV === 'production') {
        setError('This app can only be opened from Telegram');
        return;
      }

      // Get theme colors
      const themeColors = telegramApp.getTheme();
      setTheme(themeColors);
      console.log('Theme colors:', themeColors);

      // Apply theme colors to CSS variables
      applyTheme(themeColors);

      setIsInitialized(true);
      console.log('App initialized successfully');

    } catch (err) {
      console.error('Failed to initialize app:', err);
      setError('Failed to initialize the app');
    }
  };

  const applyTheme = (themeColors) => {
    const root = document.documentElement;
    
    // Set CSS variables for theming
    root.style.setProperty('--tg-theme-bg-color', themeColors.bg_color || '#ffffff');
    root.style.setProperty('--tg-theme-text-color', themeColors.text_color || '#000000');
    root.style.setProperty('--tg-theme-hint-color', themeColors.hint_color || '#999999');
    root.style.setProperty('--tg-theme-link-color', themeColors.link_color || '#007aff');
    root.style.setProperty('--tg-theme-button-color', themeColors.button_color || '#007aff');
    root.style.setProperty('--tg-theme-button-text-color', themeColors.button_text_color || '#ffffff');
    root.style.setProperty('--tg-theme-secondary-bg-color', themeColors.secondary_bg_color || '#f1f1f1');

    // Set body background
    document.body.style.backgroundColor = themeColors.bg_color || '#ffffff';
    document.body.style.color = themeColors.text_color || '#000000';
  };

  const handleRetryInit = () => {
    setError(null);
    setIsInitialized(false);
    telegramApp.hapticFeedback('impact', 'light');
    setTimeout(initializeApp, 100);
  };

  // Loading state
  if (!isInitialized && !error) {
    return (
      <div className="app-loading">
        <div className="loading-spinner"></div>
        <p>Initializing Telegram Mini-App...</p>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="app-error">
        <div className="error-content">
          <div className="error-icon">⚠️</div>
          <h2>Initialization Failed</h2>
          <p>{error}</p>
          {process.env.NODE_ENV === 'development' && (
            <div className="dev-info">
              <p>Development mode: Some features may not work outside Telegram</p>
            </div>
          )}
          <button onClick={handleRetryInit} className="retry-button">
            Try Again
          </button>
        </div>
      </div>
    );
  }

  // Main app
  return (
    <div className="app" style={{ backgroundColor: theme.bg_color }}>
      <Header />
      <main className="app-main">
        <TagList />
      </main>
    </div>
  );
}

export default App;