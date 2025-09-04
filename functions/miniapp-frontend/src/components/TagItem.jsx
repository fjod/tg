import React from 'react';
import telegramApp from '../utils/telegram.js';

const TagItem = ({ tag, onClick }) => {
  const theme = telegramApp.getTheme();

  const handleClick = () => {
    telegramApp.hapticFeedback('selection');
    if (onClick) {
      onClick(tag);
    }
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const getColorIndicator = () => {
    if (tag.color) {
      return (
        <div 
          className="tag-color-indicator"
          style={{ backgroundColor: tag.color }}
        />
      );
    }
    return (
      <div 
        className="tag-color-indicator default"
        style={{ backgroundColor: theme.button_color }}
      />
    );
  };

  return (
    <div 
      className="tag-item"
      onClick={handleClick}
      style={{
        backgroundColor: theme.bg_color,
        borderColor: theme.hint_color + '30',
        color: theme.text_color,
      }}
    >
      <div className="tag-item-header">
        <div className="tag-name-container">
          {getColorIndicator()}
          <span className="tag-name" style={{ color: theme.text_color }}>
            {tag.name}
          </span>
        </div>
        <span 
          className="tag-message-count"
          style={{ 
            color: theme.hint_color,
            backgroundColor: theme.hint_color + '20'
          }}
        >
          ({tag.message_count})
        </span>
      </div>
      <div className="tag-item-footer">
        <span className="tag-created-date" style={{ color: theme.hint_color }}>
          Created: {formatDate(tag.created_at)}
        </span>
      </div>
    </div>
  );
};

export default TagItem;