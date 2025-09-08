import React, { useState, useEffect } from 'react';
import TagItem from './TagItem.jsx';
import apiService from '../services/api.js';
import telegramApp from '../utils/telegram.js';

const TagList = ({ healthCheckResult }) => {
  const [tags, setTags] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [debugInfo, setDebugInfo] = useState(null);
  const theme = telegramApp.getTheme();

  useEffect(() => {
    loadTags();
  }, []);

  const loadTags = async () => {
    try {
      setLoading(true);
      setError(null);
      setDebugInfo(null);

      // First check if health check was successful
      if (healthCheckResult && !healthCheckResult.success) {
        console.warn('Health check failed, attempting to load tags anyway...');
        setDebugInfo({
          healthCheckStatus: 'failed',
          healthCheckError: healthCheckResult.error,
          attemptingTagsAnyway: true
        });
      } else if (healthCheckResult && healthCheckResult.success) {
        console.log('Health check passed, loading tags...');
        setDebugInfo({
          healthCheckStatus: 'passed',
          apiConnectivity: 'confirmed'
        });
      }

      // Check authentication data
      const authData = telegramApp.getAuthData();
      const user = telegramApp.getUser();
      
      console.log('Authentication debug:', {
        hasAuthData: !!authData,
        authDataLength: authData?.length || 0,
        hasUser: !!user,
        userId: user?.id,
        isInTelegram: telegramApp.isInTelegram()
      });

      setDebugInfo(prev => ({
        ...prev,
        hasAuthData: !!authData,
        authDataLength: authData?.length || 0,
        hasUser: !!user,
        userId: user?.id,
        isInTelegram: telegramApp.isInTelegram(),
        apiUrl: apiService.getBaseURL()
      }));
      
      console.log('Loading tags...');
      const userTags = await apiService.getUserTags();
      
      console.log('Tags loaded successfully:', userTags);
      setTags(userTags);
      
      setDebugInfo(prev => ({
        ...prev,
        tagsLoaded: true,
        tagCount: userTags?.length || 0
      }));
      
      // Show success haptic feedback
      telegramApp.hapticFeedback('notification', 'success');
      
    } catch (err) {
      console.error('Failed to load tags:', err);
      setError(err.message);
      setDebugInfo(prev => ({
        ...prev,
        tagsLoaded: false,
        tagError: err.message
      }));
      telegramApp.hapticFeedback('notification', 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleTagClick = (tag) => {
    console.log('Tag clicked:', tag);
    telegramApp.showAlert(`Tag: ${tag.name}\nMessages: ${tag.message_count}`);
  };

  const handleRetry = () => {
    telegramApp.hapticFeedback('impact', 'light');
    loadTags();
  };

  if (loading) {
    return (
      <div className="tag-list-container">
        <div 
          className="loading-state"
          style={{ color: theme.hint_color }}
        >
          <div className="loading-spinner"></div>
          <p>Loading your tags...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="tag-list-container">
        <div 
          className="error-state"
          style={{ 
            backgroundColor: theme.bg_color,
            color: theme.text_color,
            borderColor: '#ff6b6b'
          }}
        >
          <div className="error-icon">⚠️</div>
          <h3 style={{ color: '#ff6b6b' }}>Failed to Load Tags</h3>
          <p style={{ color: theme.hint_color }}>{error}</p>
          <button 
            className="retry-button"
            onClick={handleRetry}
            style={{
              backgroundColor: theme.button_color,
              color: theme.button_text_color,
            }}
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  if (!tags || tags.length === 0) {
    return (
      <div className="tag-list-container">
        <div 
          className="empty-state"
          style={{ 
            backgroundColor: theme.bg_color,
            color: theme.text_color 
          }}
        >
          <div className="empty-icon">🏷️</div>
          <h3>No Tags Yet</h3>
          <p style={{ color: theme.hint_color }}>
            Start organizing your messages by creating tags in the bot chat!
          </p>
          <div className="empty-instructions">
            <p style={{ color: theme.hint_color }}>
              💡 Send any message to the bot and select or create a tag
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="tag-list-container">
      {/* Debug Information Display */}
      {debugInfo && (
        <div 
          className="debug-info"
          style={{
            backgroundColor: theme.secondary_bg_color,
            color: theme.text_color,
            margin: '10px',
            padding: '8px',
            borderRadius: '6px',
            fontSize: '11px',
            border: `1px solid ${theme.hint_color}`,
            fontFamily: 'monospace'
          }}
        >
          <div style={{ fontWeight: 'bold', marginBottom: '5px', color: theme.link_color }}>
            🔍 Debug Information
          </div>
          {debugInfo.healthCheckStatus && (
            <div>Health Check: {debugInfo.healthCheckStatus} {debugInfo.healthCheckStatus === 'passed' ? '✅' : '❌'}</div>
          )}
          <div>In Telegram: {debugInfo.isInTelegram ? 'Yes ✅' : 'No ❌'}</div>
          <div>Has Auth Data: {debugInfo.hasAuthData ? 'Yes ✅' : 'No ❌'} {debugInfo.hasAuthData && `(${debugInfo.authDataLength} chars)`}</div>
          <div>Has User: {debugInfo.hasUser ? 'Yes ✅' : 'No ❌'} {debugInfo.hasUser && `(ID: ${debugInfo.userId})`}</div>
          <div>API URL: {debugInfo.apiUrl}</div>
          {debugInfo.tagsLoaded !== undefined && (
            <div>Tags Loaded: {debugInfo.tagsLoaded ? `Yes ✅ (${debugInfo.tagCount})` : `No ❌ - ${debugInfo.tagError}`}</div>
          )}
        </div>
      )}

      <div className="tag-list-header">
        <p 
          className="tag-count"
          style={{ color: theme.hint_color }}
        >
          {tags.length} tag{tags.length !== 1 ? 's' : ''} found
        </p>
      </div>
      
      <div className="tag-list">
        {tags.map((tag) => (
          <TagItem 
            key={tag.id} 
            tag={tag} 
            onClick={handleTagClick}
          />
        ))}
      </div>
      
      <div className="tag-list-footer">
        <button 
          className="refresh-button"
          onClick={handleRetry}
          style={{
            backgroundColor: 'transparent',
            color: theme.link_color,
            borderColor: theme.link_color,
          }}
        >
          🔄 Refresh
        </button>
      </div>
    </div>
  );
};

export default TagList;