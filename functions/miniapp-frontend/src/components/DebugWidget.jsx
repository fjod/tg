import React, { useState } from 'react';
import { useError } from '../contexts/ErrorContext.jsx';
import telegramApp from '../utils/telegram.js';
import apiService from '../services/api.js';

const DebugWidget = () => {
  const { 
    hasErrors, 
    getCurrentError, 
    errorHistory, 
    healthStatus,
    errors
  } = useError();
  
  const [isExpanded, setIsExpanded] = useState(false);
  const theme = telegramApp.getTheme();

  // Only show when there are errors
  if (!hasErrors) {
    return null;
  }

  const currentError = getCurrentError();
  
  const toggleExpanded = () => {
    setIsExpanded(!isExpanded);
    telegramApp.hapticFeedback('impact', 'light');
  };

  const formatTimestamp = (timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString() + '.' + date.getMilliseconds().toString().padStart(3, '0');
  };

  const getAuthStatus = () => {
    const authData = telegramApp.getAuthData();
    const user = telegramApp.getUser();
    
    return {
      inTelegram: telegramApp.isInTelegram(),
      hasAuthData: !!authData,
      authDataLength: authData?.length || 0,
      hasUser: !!user,
      userId: user?.id
    };
  };

  const authStatus = getAuthStatus();

  return (
    <div 
      className="debug-widget"
      style={{
        backgroundColor: theme.secondary_bg_color,
        color: theme.text_color,
        margin: '16px',
        borderRadius: '12px',
        border: `2px solid ${theme.hint_color}`,
        overflow: 'hidden',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
      }}
    >
      {/* Header with toggle */}
      <div 
        onClick={toggleExpanded}
        style={{
          padding: '16px',
          cursor: 'pointer',
          borderBottom: isExpanded ? `1px solid ${theme.hint_color}` : 'none',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          backgroundColor: theme.bg_color
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span style={{ fontSize: '16px' }}>🔍</span>
          <span style={{ fontWeight: 'bold', color: '#ff6b6b' }}>
            Debug Information
          </span>
        </div>
        <span style={{ 
          fontSize: '18px', 
          transform: isExpanded ? 'rotate(180deg)' : 'rotate(0deg)',
          transition: 'transform 0.2s ease'
        }}>
          ▼
        </span>
      </div>

      {/* Current Error Summary */}
      <div style={{ padding: '16px', backgroundColor: theme.bg_color }}>
        <div style={{ marginBottom: '12px' }}>
          <div style={{ 
            color: '#ff6b6b', 
            fontWeight: 'bold', 
            marginBottom: '4px' 
          }}>
            Current Error:
          </div>
          <div style={{ 
            fontFamily: 'monospace', 
            fontSize: '13px',
            backgroundColor: theme.secondary_bg_color,
            padding: '8px',
            borderRadius: '6px',
            border: '1px solid #ff6b6b'
          }}>
            <div><strong>Type:</strong> {currentError?.type}</div>
            <div><strong>Message:</strong> {currentError?.message}</div>
            <div><strong>Time:</strong> {formatTimestamp(currentError?.timestamp)}</div>
          </div>
        </div>
      </div>

      {/* Expandable Details */}
      {isExpanded && (
        <div style={{ padding: '16px', fontSize: '12px', fontFamily: 'monospace' }}>
          
          {/* API Status */}
          <div style={{ marginBottom: '16px' }}>
            <div style={{ 
              fontWeight: 'bold', 
              marginBottom: '8px', 
              color: theme.link_color 
            }}>
              📡 API Status
            </div>
            <div style={{ 
              backgroundColor: theme.secondary_bg_color, 
              padding: '8px', 
              borderRadius: '6px' 
            }}>
              <div>Base URL: {apiService.getBaseURL()}</div>
              <div>Health API: {healthStatus.api || 'unknown'}</div>
              <div>Connectivity: {healthStatus.connectivity || 'unknown'}</div>
              <div>
                Last Success: {healthStatus.lastSuccessfulCall?.timestamp 
                  ? formatTimestamp(healthStatus.lastSuccessfulCall.timestamp)
                  : 'Never'
                }
              </div>
            </div>
          </div>

          {/* Authentication Status */}
          <div style={{ marginBottom: '16px' }}>
            <div style={{ 
              fontWeight: 'bold', 
              marginBottom: '8px', 
              color: theme.link_color 
            }}>
              🔐 Authentication Status
            </div>
            <div style={{ 
              backgroundColor: theme.secondary_bg_color, 
              padding: '8px', 
              borderRadius: '6px' 
            }}>
              <div>In Telegram: {authStatus.inTelegram ? '✅ Yes' : '❌ No'}</div>
              <div>Has Auth Data: {authStatus.hasAuthData ? '✅ Yes' : '❌ No'} 
                {authStatus.hasAuthData && ` (${authStatus.authDataLength} chars)`}
              </div>
              <div>Has User: {authStatus.hasUser ? '✅ Yes' : '❌ No'}
                {authStatus.hasUser && ` (ID: ${authStatus.userId})`}
              </div>
            </div>
          </div>

          {/* Network Status */}
          <div style={{ marginBottom: '16px' }}>
            <div style={{ 
              fontWeight: 'bold', 
              marginBottom: '8px', 
              color: theme.link_color 
            }}>
              🌐 Network Status
            </div>
            <div style={{ 
              backgroundColor: theme.secondary_bg_color, 
              padding: '8px', 
              borderRadius: '6px' 
            }}>
              <div>Browser Online: {navigator.onLine ? '✅ Yes' : '❌ No'}</div>
              <div>Connection Type: {navigator.connection?.effectiveType || 'unknown'}</div>
              <div>User Agent: {navigator.userAgent.substring(0, 60)}...</div>
            </div>
          </div>

          {/* Error Details */}
          {currentError?.details && Object.keys(currentError.details).length > 0 && (
            <div style={{ marginBottom: '16px' }}>
              <div style={{ 
                fontWeight: 'bold', 
                marginBottom: '8px', 
                color: theme.link_color 
              }}>
                ⚠️ Error Details
              </div>
              <div style={{ 
                backgroundColor: theme.secondary_bg_color, 
                padding: '8px', 
                borderRadius: '6px',
                maxHeight: '120px',
                overflowY: 'auto'
              }}>
                <pre style={{ 
                  margin: 0, 
                  fontSize: '11px', 
                  whiteSpace: 'pre-wrap' 
                }}>
                  {JSON.stringify(currentError.details, null, 2)}
                </pre>
              </div>
            </div>
          )}

          {/* Error History */}
          {errorHistory.length > 1 && (
            <div>
              <div style={{ 
                fontWeight: 'bold', 
                marginBottom: '8px', 
                color: theme.link_color 
              }}>
                📝 Recent Errors ({errorHistory.length})
              </div>
              <div style={{ 
                backgroundColor: theme.secondary_bg_color, 
                padding: '8px', 
                borderRadius: '6px',
                maxHeight: '120px',
                overflowY: 'auto'
              }}>
                {errorHistory.slice(1, 4).map((error, index) => (
                  <div key={error.id} style={{ 
                    marginBottom: '6px',
                    paddingBottom: '6px',
                    borderBottom: index < 2 ? `1px solid ${theme.hint_color}` : 'none'
                  }}>
                    <div><strong>{error.type}:</strong> {error.message}</div>
                    <div style={{ color: theme.hint_color }}>
                      {formatTimestamp(error.timestamp)}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Development Mode Extra Info */}
          {process.env.NODE_ENV === 'development' && (
            <div style={{ 
              marginTop: '16px',
              padding: '8px',
              backgroundColor: '#fff3cd',
              color: '#856404',
              borderRadius: '6px',
              fontSize: '11px'
            }}>
              <strong>Development Mode:</strong> Additional debug information available in console
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default DebugWidget;