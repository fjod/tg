import React, { useEffect, useState } from 'react';
import Header from './components/Header.jsx';
import TagList from './components/TagList.jsx';
import MessageList from './components/MessageList.jsx';
import HealthCheckWidget from './components/HealthCheckWidget.jsx';
import DebugWidget from './components/DebugWidget.jsx';
import { ErrorProvider, useError } from './contexts/ErrorContext.jsx';
import { NavigationProvider, useNavigation } from './contexts/NavigationContext.jsx';
import telegramApp from './utils/telegram.js';
import apiService from './services/api.js';
import './styles.css';

function AppContent() {
  const [isInitialized, setIsInitialized] = useState(false);
  const [theme, setTheme] = useState({});
  const { addError, clearAllErrors, updateHealthStatus, markSuccessfulCall } = useError();
  const { currentView, isInTagsView, isInMessagesView } = useNavigation();

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
        addError('general', 'This app should be opened from Telegram');
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
        addError('auth', 'This app can only be opened from Telegram');
        return;
      }

      // Get theme colors
      const themeColors = telegramApp.getTheme();
      setTheme(themeColors);
      console.log('Theme colors:', themeColors);

      // Apply theme colors to CSS variables
      applyTheme(themeColors);

      // Perform health check to test API connectivity
      console.log('Performing health check...');
      try {
        const healthUrl = `${apiService.getBaseURL()}/api/health`;
        console.log('Health check URL:', healthUrl);
        
        const response = await fetch(healthUrl, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          timeout: 10000,
        });
        
        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        const healthResponse = await response.json();
        console.log('Health check successful:', healthResponse);
        
        // Mark successful API call
        markSuccessfulCall('/api/health');
        updateHealthStatus({
          api: 'healthy',
          connectivity: 'online'
        });
        
      } catch (healthError) {
        console.error('Health check failed:', healthError);
        addError('api', healthError, {
          endpoint: '/api/health',
          url: `${apiService.getBaseURL()}/api/health`,
          method: 'GET'
        });
        
        updateHealthStatus({
          api: 'error',
          connectivity: 'error'
        });
      }

      setIsInitialized(true);
      console.log('App initialized successfully');

    } catch (err) {
      console.error('Failed to initialize app:', err);
      addError('general', 'Failed to initialize the app', { originalError: err.message });
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
    clearAllErrors();
    setIsInitialized(false);
    telegramApp.hapticFeedback('impact', 'light');
    setTimeout(initializeApp, 100);
  };

  const { hasErrors, getCurrentError } = useError();
  const currentError = getCurrentError();

  // Loading state
  if (!isInitialized && !hasErrors) {
    return (
      <div className="app-loading">
        <div className="loading-spinner"></div>
        <p>Initializing Telegram Mini-App...</p>
      </div>
    );
  }

  // Error state during initialization
  if (!isInitialized && hasErrors) {
    return (
      <div className="app-error">
        <div className="error-content">
          <div className="error-icon">⚠️</div>
          <h2>Initialization Failed</h2>
          <p>{currentError?.message || 'Unknown error occurred'}</p>
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

  // Render current view based on navigation state
  const renderCurrentView = () => {
    switch (currentView) {
      case 'messages':
        return <MessageList />;
      case 'loading':
        return (
          <div className="app-loading">
            <div className="loading-spinner"></div>
            <p>Loading...</p>
          </div>
        );
      case 'tags':
      default:
        return <TagList />;
    }
  };

  // Main app
  return (
    <div className="app" style={{ backgroundColor: theme.bg_color }}>
      <Header />
      <main className="app-main">
        {/* Conditional error widgets - only show when there are errors */}
        <HealthCheckWidget />
        <DebugWidget />
        
        {/* Render current view */}
        {renderCurrentView()}
      </main>
    </div>
  );
}

// Wrapper component with providers
function App() {
  return (
    <ErrorProvider>
      <NavigationProvider>
        <AppContent />
      </NavigationProvider>
    </ErrorProvider>
  );
}

export default App;