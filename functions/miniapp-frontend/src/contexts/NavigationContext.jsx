import React, { createContext, useContext, useState, useCallback } from 'react';
import telegramApp from '../utils/telegram.js';

const NavigationContext = createContext();

export const useNavigation = () => {
  const context = useContext(NavigationContext);
  if (!context) {
    throw new Error('useNavigation must be used within a NavigationProvider');
  }
  return context;
};

export const NavigationProvider = ({ children }) => {
  const [currentView, setCurrentView] = useState('tags'); // 'tags' | 'messages' | 'loading'
  const [selectedTag, setSelectedTag] = useState(null);
  const [navigationHistory, setNavigationHistory] = useState(['tags']);
  const [isLoading, setIsLoading] = useState(false);

  // Navigate to messages view for a specific tag
  const navigateToMessages = useCallback((tag) => {
    console.log('NavigationContext: navigateToMessages called with tag:', tag);
    
    try {
      // Haptic feedback for navigation
      console.log('NavigationContext: calling hapticFeedback');
      telegramApp.hapticFeedback('selection');
      
      // Update state
      console.log('NavigationContext: updating state');
      setSelectedTag(tag);
      setCurrentView('messages');
      setNavigationHistory(prev => {
        console.log('NavigationContext: updating history from:', prev);
        return [...prev, 'messages'];
      });
      
      // Update Telegram WebApp back button
      console.log('NavigationContext: setting up back button');
      if (telegramApp.tg?.BackButton) {
        telegramApp.tg.BackButton.show();
        telegramApp.tg.BackButton.onClick(() => navigateBack());
      }
      
      console.log('NavigationContext: navigateToMessages completed successfully');
    } catch (error) {
      console.error('NavigationContext: Error in navigateToMessages:', error);
      // Don't re-throw to avoid breaking the UI
    }
  }, [navigateBack]);

  // Navigate back to previous view
  const navigateBack = useCallback(() => {
    console.log('Navigating back, current history:', navigationHistory);
    
    // Haptic feedback
    telegramApp.hapticFeedback('impact', 'light');
    
    if (navigationHistory.length > 1) {
      const newHistory = [...navigationHistory];
      newHistory.pop(); // Remove current view
      const previousView = newHistory[newHistory.length - 1];
      
      setNavigationHistory(newHistory);
      setCurrentView(previousView);
      
      // Clear selected tag when going back to tags
      if (previousView === 'tags') {
        setSelectedTag(null);
        // Hide back button when at root
        if (telegramApp.tg?.BackButton) {
          telegramApp.tg.BackButton.hide();
        }
      }
    } else {
      // Already at root, just ensure we're in tags view
      setCurrentView('tags');
      setSelectedTag(null);
      if (telegramApp.tg?.BackButton) {
        telegramApp.tg.BackButton.hide();
      }
    }
  }, [navigationHistory]);

  // Reset navigation to root (tags view)
  const resetNavigation = useCallback(() => {
    console.log('Resetting navigation to tags view');
    
    setCurrentView('tags');
    setSelectedTag(null);
    setNavigationHistory(['tags']);
    setIsLoading(false);
    
    // Hide back button
    if (telegramApp.tg?.BackButton) {
      telegramApp.tg.BackButton.hide();
    }
  }, []);

  // Set loading state
  const setNavigationLoading = useCallback((loading) => {
    setIsLoading(loading);
    if (loading) {
      setCurrentView('loading');
    }
  }, []);

  // Get current view context
  const getCurrentContext = useCallback(() => {
    return {
      view: currentView,
      tag: selectedTag,
      isAtRoot: currentView === 'tags',
      canGoBack: navigationHistory.length > 1,
      isLoading
    };
  }, [currentView, selectedTag, navigationHistory, isLoading]);

  // Get breadcrumb for current location
  const getBreadcrumb = useCallback(() => {
    const breadcrumb = [];
    
    if (navigationHistory.includes('tags')) {
      breadcrumb.push({ label: 'Tags', view: 'tags' });
    }
    
    if (navigationHistory.includes('messages') && selectedTag) {
      breadcrumb.push({ 
        label: selectedTag.name, 
        view: 'messages',
        tag: selectedTag 
      });
    }
    
    return breadcrumb;
  }, [navigationHistory, selectedTag]);

  // Handle browser back button / hardware back button
  const handleBackButton = useCallback(() => {
    navigateBack();
  }, [navigateBack]);

  // Initialize back button handling
  React.useEffect(() => {
    // Handle Telegram back button
    if (telegramApp.tg?.BackButton) {
      telegramApp.tg.BackButton.onClick(handleBackButton);
    }
    
    // Handle browser back button
    const handlePopState = () => {
      handleBackButton();
    };
    
    window.addEventListener('popstate', handlePopState);
    
    return () => {
      window.removeEventListener('popstate', handlePopState);
      if (telegramApp.tg?.BackButton) {
        telegramApp.tg.BackButton.offClick(handleBackButton);
      }
    };
  }, [handleBackButton]);

  const value = {
    // State
    currentView,
    selectedTag,
    navigationHistory,
    isLoading,

    // Actions
    navigateToMessages,
    navigateBack,
    resetNavigation,
    setNavigationLoading,

    // Computed
    getCurrentContext,
    getBreadcrumb,

    // Convenience properties
    isInTagsView: currentView === 'tags',
    isInMessagesView: currentView === 'messages',
    isInLoadingView: currentView === 'loading',
    canGoBack: navigationHistory.length > 1
  };

  return (
    <NavigationContext.Provider value={value}>
      {children}
    </NavigationContext.Provider>
  );
};