import React, { useState, useEffect } from 'react';
import TagItem from './TagItem.jsx';
import { useError } from '../contexts/ErrorContext.jsx';
import { useNavigation } from '../contexts/NavigationContext.jsx';
import apiService from '../services/api.js';
import telegramApp from '../utils/telegram.js';

const TagList = () => {
  const [tags, setTags] = useState([]);
  const [loading, setLoading] = useState(true);
  const { addError, clearError, markSuccessfulCall, hasApiError } = useError();
  const { navigateToMessages } = useNavigation();
  const theme = telegramApp.getTheme();

  useEffect(() => {
    loadTags();
  }, []);

  const loadTags = async () => {
    try {
      setLoading(true);
      clearError('api');

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

      // Validate authentication before making API call
      if (!authData && telegramApp.isInTelegram()) {
        throw new Error('No authentication data available');
      }
      
      console.log('Loading tags...');
      const userTags = await apiService.getUserTags();
      
      console.log('Tags loaded successfully:', userTags);
      setTags(userTags);
      
      // Mark successful API call
      markSuccessfulCall('/api/user/tags');
      
      // Show success haptic feedback
      telegramApp.hapticFeedback('notification', 'success');
      
    } catch (err) {
      console.error('Failed to load tags:', err);
      
      // Determine error type
      const errorType = err.message?.includes('auth') ? 'auth' : 'api';
      
      addError(errorType, err, {
        endpoint: '/api/user/tags',
        method: 'GET',
        hasAuthData: !!telegramApp.getAuthData(),
        isInTelegram: telegramApp.isInTelegram()
      });
      
      telegramApp.hapticFeedback('notification', 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleTagClick = (tag) => {
    console.log('Tag clicked:', tag);
    navigateToMessages(tag);
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