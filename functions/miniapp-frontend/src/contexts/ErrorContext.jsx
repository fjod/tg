import React, { createContext, useContext, useState, useCallback } from 'react';

const ErrorContext = createContext();

export const useError = () => {
  const context = useContext(ErrorContext);
  if (!context) {
    throw new Error('useError must be used within an ErrorProvider');
  }
  return context;
};

export const ErrorProvider = ({ children }) => {
  const [errors, setErrors] = useState({
    api: null,
    network: null,
    auth: null,
    general: null
  });
  
  const [errorHistory, setErrorHistory] = useState([]);
  const [healthStatus, setHealthStatus] = useState({
    api: null,
    lastSuccessfulCall: null,
    connectivity: 'unknown'
  });

  // Add error with timestamp and type
  const addError = useCallback((type, error, details = {}) => {
    const errorEntry = {
      type,
      message: error.message || error,
      details,
      timestamp: new Date().toISOString(),
      id: Date.now()
    };

    // Update current errors
    setErrors(prev => ({
      ...prev,
      [type]: errorEntry
    }));

    // Add to history (keep last 10 errors)
    setErrorHistory(prev => [errorEntry, ...prev.slice(0, 9)]);

    console.error(`[ErrorContext] ${type} error:`, error, details);
  }, []);

  // Clear specific error type
  const clearError = useCallback((type) => {
    setErrors(prev => ({
      ...prev,
      [type]: null
    }));
  }, []);

  // Clear all errors
  const clearAllErrors = useCallback(() => {
    setErrors({
      api: null,
      network: null,
      auth: null,
      general: null
    });
  }, []);

  // Update health status
  const updateHealthStatus = useCallback((status) => {
    setHealthStatus(prev => ({
      ...prev,
      ...status,
      lastUpdated: new Date().toISOString()
    }));
  }, []);

  // Mark successful API call
  const markSuccessfulCall = useCallback((endpoint) => {
    setHealthStatus(prev => ({
      ...prev,
      api: 'healthy',
      lastSuccessfulCall: {
        endpoint,
        timestamp: new Date().toISOString()
      },
      connectivity: 'online'
    }));

    // Clear API and network errors on success
    clearError('api');
    clearError('network');
  }, [clearError]);

  // Check if any errors exist
  const hasErrors = Object.values(errors).some(error => error !== null);
  
  // Check if specific error types exist
  const hasApiError = errors.api !== null;
  const hasNetworkError = errors.network !== null;
  const hasAuthError = errors.auth !== null;
  const hasGeneralError = errors.general !== null;

  // Get current error for display
  const getCurrentError = () => {
    // Priority: auth > api > network > general
    if (errors.auth) return errors.auth;
    if (errors.api) return errors.api;
    if (errors.network) return errors.network;
    if (errors.general) return errors.general;
    return null;
  };

  // Retry helpers
  const createRetryHandler = useCallback((retryFn, errorType) => {
    return async () => {
      try {
        clearError(errorType);
        await retryFn();
      } catch (error) {
        addError(errorType, error);
      }
    };
  }, [addError, clearError]);

  const value = {
    // State
    errors,
    errorHistory,
    healthStatus,
    hasErrors,
    hasApiError,
    hasNetworkError,
    hasAuthError,
    hasGeneralError,

    // Actions
    addError,
    clearError,
    clearAllErrors,
    updateHealthStatus,
    markSuccessfulCall,
    getCurrentError,
    createRetryHandler
  };

  return (
    <ErrorContext.Provider value={value}>
      {children}
    </ErrorContext.Provider>
  );
};