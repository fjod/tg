/**
 * Utility functions for message formatting and display helpers
 */

/**
 * Get message type icon based on message type
 * @param {string} messageType - The message type (text, photo, video, etc.)
 * @returns {string} Emoji icon for the message type
 */
export const getMessageTypeIcon = (messageType) => {
  const iconMap = {
    'text': '💬',
    'photo': '📷',
    'video': '🎥',
    'document': '📄',
    'audio': '🎵',
    'voice': '🎤',
    'video_note': '📹',
    'animation': '🎬',
    'sticker': '🏷️',
    'location': '📍',
    'contact': '👤',
    'poll': '📊',
    'dice': '🎲',
    'game': '🎮',
    'invoice': '🧾'
  };
  
  return iconMap[messageType] || '📝';
};

/**
 * Get message preview text based on content
 * @param {Object} message - The message object
 * @returns {string} Preview text for the message
 */
export const getMessagePreview = (message) => {
  const { message_type, text_content, caption, file_name } = message;
  
  // For text messages, use the text content
  if (message_type === 'text' && text_content) {
    return truncateText(text_content, 100);
  }
  
  // For messages with captions, use the caption
  if (caption) {
    return truncateText(caption, 100);
  }
  
  // For files with names, show the filename
  if (file_name) {
    return file_name;
  }
  
  // For specific message types without text content
  const typePreviewMap = {
    'photo': 'Photo',
    'video': 'Video', 
    'audio': 'Audio file',
    'voice': 'Voice message',
    'video_note': 'Video note',
    'animation': 'GIF',
    'sticker': 'Sticker',
    'document': 'Document',
    'location': 'Location',
    'contact': 'Contact',
    'poll': 'Poll',
    'dice': 'Dice',
    'game': 'Game',
    'invoice': 'Invoice'
  };
  
  return typePreviewMap[message_type] || 'Message';
};

/**
 * Truncate text to specified length with ellipsis
 * @param {string} text - Text to truncate
 * @param {number} maxLength - Maximum length before truncation
 * @returns {string} Truncated text
 */
export const truncateText = (text, maxLength = 100) => {
  if (!text || text.length <= maxLength) {
    return text || '';
  }
  
  return text.substring(0, maxLength).trim() + '...';
};

/**
 * Format date for message display
 * @param {string|Date} dateString - Date to format
 * @returns {string} Formatted date string
 */
export const formatMessageDate = (dateString) => {
  const date = new Date(dateString);
  const now = new Date();
  
  // If date is invalid, return fallback
  if (isNaN(date.getTime())) {
    return 'Unknown date';
  }
  
  const diffMs = now - date;
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffMinutes = Math.floor(diffMs / (1000 * 60));
  
  // Less than 1 minute ago
  if (diffMinutes < 1) {
    return 'Just now';
  }
  
  // Less than 1 hour ago
  if (diffMinutes < 60) {
    return `${diffMinutes}m ago`;
  }
  
  // Less than 24 hours ago
  if (diffHours < 24) {
    return `${diffHours}h ago`;
  }
  
  // Less than 7 days ago
  if (diffDays < 7) {
    return `${diffDays}d ago`;
  }
  
  // More than a week ago, show formatted date
  if (diffDays < 365) {
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric'
    });
  }
  
  // More than a year ago, include year
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  });
};

/**
 * Format time for message display (just time part)
 * @param {string|Date} dateString - Date to format
 * @returns {string} Formatted time string
 */
export const formatMessageTime = (dateString) => {
  const date = new Date(dateString);
  
  if (isNaN(date.getTime())) {
    return '';
  }
  
  return date.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  });
};

/**
 * Get message metadata display text
 * @param {Object} message - The message object
 * @returns {string} Metadata text (forwarded info, etc.)
 */
export const getMessageMetadata = (message) => {
  const metadata = [];
  
  if (message.forwarded_from) {
    metadata.push(`Forwarded from ${message.forwarded_from}`);
  }
  
  if (message.urls && message.urls.length > 0) {
    metadata.push(`${message.urls.length} link${message.urls.length > 1 ? 's' : ''}`);
  }
  
  if (message.hashtags && message.hashtags.length > 0) {
    metadata.push(`${message.hashtags.length} hashtag${message.hashtags.length > 1 ? 's' : ''}`);
  }
  
  return metadata.join(' • ');
};

/**
 * Check if message has media content
 * @param {Object} message - The message object
 * @returns {boolean} True if message contains media
 */
export const isMediaMessage = (message) => {
  const mediaTypes = ['photo', 'video', 'audio', 'voice', 'video_note', 'animation', 'sticker'];
  return mediaTypes.includes(message.message_type);
};

/**
 * Check if message has file attachment
 * @param {Object} message - The message object
 * @returns {boolean} True if message has file
 */
export const hasFileAttachment = (message) => {
  return !!(message.file_name || message.file_id);
};

/**
 * Get file size display text
 * @param {number} sizeInBytes - File size in bytes
 * @returns {string} Formatted file size
 */
export const formatFileSize = (sizeInBytes) => {
  if (!sizeInBytes || sizeInBytes === 0) {
    return '';
  }
  
  const units = ['B', 'KB', 'MB', 'GB'];
  let size = sizeInBytes;
  let unitIndex = 0;
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }
  
  return `${size.toFixed(unitIndex > 0 ? 1 : 0)} ${units[unitIndex]}`;
};

/**
 * Generate Telegram message URL for redirection
 * @param {Object} message - The message object with telegram_message_id
 * @param {number} userId - User's Telegram ID
 * @returns {string} Telegram message URL
 */
export const generateTelegramMessageUrl = (message, userId) => {
  if (!message.telegram_message_id || !userId) {
    return null;
  }
  
  // For private chats, use direct message link format
  // Note: This may need adjustment based on actual Telegram chat structure
  return `https://t.me/c/${userId}/${message.telegram_message_id}`;
};