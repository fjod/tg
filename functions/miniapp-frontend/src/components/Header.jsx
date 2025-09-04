import React from 'react';
import telegramApp from '../utils/telegram.js';

const Header = () => {
  const user = telegramApp.getUser();
  const theme = telegramApp.getTheme();

  return (
    <header 
      className="header"
      style={{ 
        backgroundColor: theme.bg_color,
        borderBottomColor: theme.hint_color + '40', // Add opacity
        color: theme.text_color 
      }}
    >
      <div className="header-content">
        <div className="header-title">
          <h1 style={{ color: theme.text_color }}>🏷️ My Tags</h1>
          {user && (
            <p className="header-subtitle" style={{ color: theme.hint_color }}>
              Welcome, {user.first_name}
            </p>
          )}
        </div>
        <div className="header-actions">
          {/* Future: Add settings or menu button */}
        </div>
      </div>
    </header>
  );
};

export default Header;