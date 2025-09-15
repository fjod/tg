import React from 'react';
import { useError } from '../contexts/ErrorContext.jsx';
import telegramApp from '../utils/telegram.js';
import {
  getMessageTypeIcon,
  getMessagePreview,
  formatMessageDate,
  formatMessageTime,
  getMessageMetadata,
  isMediaMessage,
  hasFileAttachment,
  formatFileSize
} from '../utils/messageHelpers.js';

const MessageItem = ({ message, onClick }) => {
  const theme = telegramApp.getTheme();
  const { addError } = useError();

  const handleClick = async () => {
    try {
      telegramApp.hapticFeedback('selection');
      
      if (onClick) {
        await onClick(message);
      } else {
        // Default behavior: try to redirect to Telegram
        redirectToTelegram(message);
      }
    } catch (error) {
      addError('general', error, { 
        messageId: message.id,
        action: 'message_click'
      });
    }
  };

  const redirectToTelegram = (message) => {
    try {
      const success = telegramApp.redirectToMessage(message);
      
      if (!success) {
        throw new Error('Failed to redirect to Telegram message');
      }
    } catch (error) {
      addError('general', error, { 
        messageId: message.id,
        telegramMessageId: message.telegram_message_id,
        action: 'telegram_redirect'
      });
    }
  };

  const messageIcon = getMessageTypeIcon(message.message_type);
  const preview = getMessagePreview(message);
  const metadata = getMessageMetadata(message);
  const isMedia = isMediaMessage(message);
  const hasFile = hasFileAttachment(message);

  return (
    <div 
      className="message-item"
      onClick={handleClick}
      style={{
        backgroundColor: theme.bg_color,
        borderColor: theme.hint_color + '20',
        color: theme.text_color,
        cursor: 'pointer',
        padding: '16px',
        marginBottom: '8px',
        borderRadius: '12px',
        border: '1px solid',
        transition: 'all 0.2s ease',
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)'
      }}
      onMouseDown={(e) => {
        e.currentTarget.style.backgroundColor = theme.secondary_bg_color;
      }}
      onMouseUp={(e) => {
        e.currentTarget.style.backgroundColor = theme.bg_color;
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.backgroundColor = theme.bg_color;
      }}
    >
      {/* Header with type icon and date */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        marginBottom: '8px'
      }}>
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: '8px'
        }}>
          <span style={{ fontSize: '20px' }}>{messageIcon}</span>
          <span style={{
            color: theme.hint_color,
            fontSize: '12px',
            textTransform: 'capitalize'
          }}>
            {message.message_type.replace('_', ' ')}
          </span>
        </div>
        
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: '6px',
          fontSize: '12px',
          color: theme.hint_color
        }}>
          <span>{formatMessageTime(message.created_at)}</span>
          <span>•</span>
          <span>{formatMessageDate(message.created_at)}</span>
        </div>
      </div>

      {/* Message preview content */}
      <div style={{
        marginBottom: metadata ? '8px' : '0'
      }}>
        <div style={{
          fontSize: '14px',
          lineHeight: '1.4',
          color: theme.text_color,
          wordBreak: 'break-word'
        }}>
          {preview}
        </div>

        {/* File size info for files */}
        {hasFile && message.file_size && (
          <div style={{
            fontSize: '12px',
            color: theme.hint_color,
            marginTop: '4px'
          }}>
            {formatFileSize(message.file_size)}
          </div>
        )}
      </div>

      {/* Metadata footer */}
      {metadata && (
        <div style={{
          fontSize: '12px',
          color: theme.hint_color,
          fontStyle: 'italic'
        }}>
          {metadata}
        </div>
      )}

      {/* URLs preview (if any) */}
      {message.urls && message.urls.length > 0 && (
        <div style={{
          marginTop: '8px',
          fontSize: '12px'
        }}>
          {message.urls.slice(0, 2).map((url, index) => (
            <div key={index} style={{
              color: theme.link_color,
              textDecoration: 'none',
              backgroundColor: theme.secondary_bg_color,
              padding: '4px 8px',
              borderRadius: '6px',
              marginBottom: '2px',
              wordBreak: 'break-all'
            }}>
              🔗 {url.length > 50 ? url.substring(0, 47) + '...' : url}
            </div>
          ))}
          {message.urls.length > 2 && (
            <div style={{
              color: theme.hint_color,
              fontSize: '11px',
              marginTop: '2px'
            }}>
              +{message.urls.length - 2} more link{message.urls.length - 2 > 1 ? 's' : ''}
            </div>
          )}
        </div>
      )}

      {/* Hashtags preview (if any) */}
      {message.hashtags && message.hashtags.length > 0 && (
        <div style={{
          marginTop: '6px',
          display: 'flex',
          flexWrap: 'wrap',
          gap: '4px'
        }}>
          {message.hashtags.slice(0, 5).map((hashtag, index) => (
            <span key={index} style={{
              fontSize: '11px',
              color: theme.link_color,
              backgroundColor: theme.link_color + '15',
              padding: '2px 6px',
              borderRadius: '10px'
            }}>
              #{hashtag}
            </span>
          ))}
          {message.hashtags.length > 5 && (
            <span style={{
              fontSize: '11px',
              color: theme.hint_color,
              padding: '2px 6px'
            }}>
              +{message.hashtags.length - 5}
            </span>
          )}
        </div>
      )}

      {/* Media indicator */}
      {isMedia && (
        <div style={{
          position: 'absolute',
          top: '8px',
          right: '8px',
          backgroundColor: theme.button_color + '20',
          borderRadius: '50%',
          width: '6px',
          height: '6px'
        }} />
      )}
    </div>
  );
};

export default MessageItem;