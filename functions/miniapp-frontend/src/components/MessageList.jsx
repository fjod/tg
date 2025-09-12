import React, { useState, useEffect } from 'react';
import MessageItem from './MessageItem.jsx';
import { useNavigation } from '../contexts/NavigationContext.jsx';
import { useError } from '../contexts/ErrorContext.jsx';
import apiService from '../services/api.js';
import telegramApp from '../utils/telegram.js';

const MessageList = () => {
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(true);
  const { selectedTag, navigateBack } = useNavigation();
  const { addError, clearError, hasApiError } = useError();
  const theme = telegramApp.getTheme();

  useEffect(() => {
    if (selectedTag) {
      loadMessages(selectedTag.id);
    }
  }, [selectedTag]);

  const loadMessages = async (tagId) => {
    try {
      setLoading(true);
      clearError('api');
      
      console.log('Loading messages for tag:', tagId);
      
      const tagMessages = await apiService.getTagMessages(tagId);
      setMessages(tagMessages);
      console.log('Messages loaded successfully:', tagMessages);
      
    } catch (error) {
      console.error('Failed to load messages :', error);
      addError('api', error, {
        endpoint: `/api/user/tags/${tagId}/messages`,
        tagId: tagId
      });
    } finally {
      setLoading(false);
    }
  };

  const handleRetry = () => {
    if (selectedTag) {
      telegramApp.hapticFeedback('impact', 'light');
      loadMessages(selectedTag.id);
    }
  };

  const handleBackClick = () => {
    telegramApp.hapticFeedback('impact', 'light');
    navigateBack();
  };

  const handleMessageClick = async (message) => {
    try {
      // Default behavior - redirect to Telegram
      // MessageItem component handles the actual redirection
      console.log('Message clicked:', message);
    } catch (error) {
      console.error('Error handling message click:', error);
    }
  };


  if (!selectedTag) {
    return (
      <div style={{ 
        padding: '20px', 
        textAlign: 'center',
        color: theme.hint_color 
      }}>
        No tag selected
      </div>
    );
  }

  return (
    <div className="message-list-container" style={{ height: '100%' }}>
      {/* Header with back button and tag info */}
      <div 
        className="message-list-header"
        style={{
          position: 'sticky',
          top: 0,
          backgroundColor: theme.bg_color,
          borderBottom: `1px solid ${theme.hint_color}20`,
          padding: '16px',
          zIndex: 10
        }}
      >
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: '12px'
        }}>
          {/* Back button */}
          <button
            onClick={handleBackClick}
            style={{
              backgroundColor: 'transparent',
              border: 'none',
              color: theme.link_color,
              fontSize: '20px',
              cursor: 'pointer',
              padding: '4px',
              borderRadius: '6px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}
          >
            ←
          </button>
          
          {/* Tag info */}
          <div style={{ flex: 1 }}>
            <h2 style={{
              margin: 0,
              fontSize: '18px',
              color: theme.text_color,
              fontWeight: '600'
            }}>
              {selectedTag.name}
            </h2>
            <p style={{
              margin: 0,
              fontSize: '12px',
              color: theme.hint_color
            }}>
              {loading ? 'Loading messages...' : 
               messages.length === 0 ? 'No messages' :
               `${messages.length} message${messages.length !== 1 ? 's' : ''}`}
            </p>
          </div>

          {/* Refresh button */}
          <button
            onClick={handleRetry}
            style={{
              backgroundColor: 'transparent',
              border: `1px solid ${theme.hint_color}40`,
              color: theme.link_color,
              fontSize: '14px',
              cursor: 'pointer',
              padding: '6px 12px',
              borderRadius: '16px',
              display: 'flex',
              alignItems: 'center',
              gap: '4px'
            }}
            disabled={loading}
          >
            🔄 {loading ? 'Loading...' : 'Refresh'}
          </button>
        </div>
      </div>

      {/* Content area */}
      <div style={{ 
        padding: '16px',
        paddingTop: '8px'
      }}>
        {loading ? (
          /* Loading state */
          <div style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '40px 20px',
            color: theme.hint_color
          }}>
            <div className="loading-spinner" style={{
              width: '32px',
              height: '32px',
              marginBottom: '16px'
            }} />
            <p>Loading messages...</p>
          </div>
        ) : hasApiError ? (
          /* Error state */
          <div style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '40px 20px',
            textAlign: 'center'
          }}>
            <div style={{ fontSize: '48px', marginBottom: '16px' }}>⚠️</div>
            <h3 style={{ 
              color: '#ff6b6b',
              marginBottom: '8px',
              fontSize: '16px'
            }}>
              Failed to Load Messages
            </h3>
            <p style={{ 
              color: theme.hint_color,
              marginBottom: '20px',
              fontSize: '14px'
            }}>
              Unable to fetch messages for this tag. Please check your connection and try again.
            </p>
            <button 
              onClick={handleRetry}
              style={{
                backgroundColor: theme.button_color,
                color: theme.button_text_color,
                border: 'none',
                padding: '12px 24px',
                borderRadius: '8px',
                fontSize: '14px',
                cursor: 'pointer'
              }}
            >
              Try Again
            </button>
          </div>
        ) : messages.length === 0 ? (
          /* Empty state */
          <div style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '60px 20px',
            textAlign: 'center'
          }}>
            <div style={{ fontSize: '48px', marginBottom: '16px' }}>📭</div>
            <h3 style={{ 
              color: theme.text_color,
              marginBottom: '8px',
              fontSize: '16px'
            }}>
              No Messages Found
            </h3>
            <p style={{ 
              color: theme.hint_color,
              fontSize: '14px'
            }}>
              This tag doesn't have any messages yet. Messages you tag with "{selectedTag.name}" will appear here.
            </p>
          </div>
        ) : (
          /* Messages list */
          <div className="messages-list">
            {messages.map((message) => (
              <MessageItem
                key={message.id}
                message={message}
                onClick={handleMessageClick}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default MessageList;