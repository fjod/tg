import React from 'react';
import { useError } from '../contexts/ErrorContext.jsx';
import telegramApp from '../utils/telegram.js';

const HealthCheckWidget = () => {
  const { 
    hasApiError, 
    hasNetworkError, 
    healthStatus, 
    createRetryHandler,
    clearAllErrors
  } = useError();

  const theme = telegramApp.getTheme();

  // Only show when there are API or network errors
  if (!hasApiError && !hasNetworkError) {
    return null;
  }

  const handleRetry = () => {
    telegramApp.hapticFeedback('impact', 'light');
    clearAllErrors();
    // Trigger a page reload to retry initialization
    window.location.reload();
  };

  const getConnectionStatus = () => {
    if (hasNetworkError) return 'offline';
    if (hasApiError) return 'api-error';
    return healthStatus.connectivity || 'unknown';
  };

  const getStatusIcon = () => {
    const status = getConnectionStatus();
    switch (status) {
      case 'offline': return '🔴';
      case 'api-error': return '⚠️';
      case 'online': return '🟢';
      default: return '🟡';
    }
  };

  const getStatusText = () => {
    const status = getConnectionStatus();
    switch (status) {
      case 'offline': return 'Connection Lost';
      case 'api-error': return 'API Error';
      case 'online': return 'Connected';
      default: return 'Unknown Status';
    }
  };

  const formatLastSuccess = () => {
    if (!healthStatus.lastSuccessfulCall) return 'Never';
    
    const time = new Date(healthStatus.lastSuccessfulCall.timestamp);
    const now = new Date();
    const diffMs = now - time;
    const diffSecs = Math.floor(diffMs / 1000);
    const diffMins = Math.floor(diffSecs / 60);
    
    if (diffSecs < 60) return `${diffSecs}s ago`;
    if (diffMins < 60) return `${diffMins}m ago`;
    return time.toLocaleTimeString();
  };

  return (
    <div 
      className="health-check-widget"
      style={{
        backgroundColor: theme.secondary_bg_color,
        color: theme.text_color,
        margin: '16px',
        padding: '16px',
        borderRadius: '12px',
        border: '2px solid #ff6b6b',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
      }}
    >
      <div style={{ 
        display: 'flex', 
        alignItems: 'center', 
        marginBottom: '12px',
        fontWeight: 'bold'
      }}>
        <span style={{ fontSize: '18px', marginRight: '8px' }}>
          {getStatusIcon()}
        </span>
        <span style={{ color: '#ff6b6b' }}>
          {getStatusText()}
        </span>
      </div>

      <div style={{ marginBottom: '12px', fontSize: '14px' }}>
        <div style={{ marginBottom: '6px' }}>
          <span style={{ color: theme.hint_color }}>Last Success: </span>
          <span style={{ fontWeight: '500' }}>
            {formatLastSuccess()}
          </span>
        </div>
        
        {healthStatus.lastSuccessfulCall?.endpoint && (
          <div style={{ marginBottom: '6px' }}>
            <span style={{ color: theme.hint_color }}>Endpoint: </span>
            <span style={{ 
              fontFamily: 'monospace', 
              fontSize: '12px',
              backgroundColor: theme.bg_color,
              padding: '2px 6px',
              borderRadius: '4px'
            }}>
              {healthStatus.lastSuccessfulCall.endpoint}
            </span>
          </div>
        )}

        <div>
          <span style={{ color: theme.hint_color }}>Network: </span>
          <span style={{ 
            color: navigator.onLine ? '#4CAF50' : '#ff6b6b',
            fontWeight: '500'
          }}>
            {navigator.onLine ? 'Online' : 'Offline'}
          </span>
        </div>
      </div>

      <button
        onClick={handleRetry}
        style={{
          backgroundColor: theme.button_color,
          color: theme.button_text_color,
          border: 'none',
          padding: '12px 24px',
          borderRadius: '8px',
          fontSize: '14px',
          fontWeight: '500',
          cursor: 'pointer',
          width: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          gap: '8px'
        }}
      >
        🔄 Retry Connection
      </button>

      <div style={{ 
        marginTop: '8px', 
        fontSize: '12px', 
        color: theme.hint_color,
        textAlign: 'center'
      }}>
        Check your internet connection and try again
      </div>
    </div>
  );
};

export default HealthCheckWidget;