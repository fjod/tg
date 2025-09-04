# Telegram Content Organizer Mini-App

A React-based Telegram Mini-App that allows users to view and manage their message tags.

## Features

- 🏷️ View all user tags with message counts
- 🎨 Telegram theme integration (light/dark mode support)
- 📱 Mobile-first responsive design  
- ⚡ Fast loading with optimized bundle
- 🔐 Secure authentication via Telegram WebApp
- 📊 Loading states and error handling
- 🎯 Empty state guidance for new users

## Tech Stack

- **Frontend**: React 18 with functional components and hooks
- **Styling**: CSS with CSS variables for theming  
- **Authentication**: Telegram WebApp SDK
- **API**: REST API integration with existing Go backend
- **Build**: React Scripts (Create React App)

## Project Structure

```
miniapp-frontend/
├── public/
│   ├── index.html          # Main HTML with Telegram WebApp SDK
│   └── manifest.json       # PWA manifest
├── src/
│   ├── components/
│   │   ├── Header.jsx      # App header with user info
│   │   ├── TagList.jsx     # Main tag list with states  
│   │   └── TagItem.jsx     # Individual tag component
│   ├── services/
│   │   └── api.js          # API service layer
│   ├── utils/
│   │   └── telegram.js     # Telegram WebApp utilities
│   ├── App.jsx             # Main app component
│   ├── index.js            # React entry point
│   └── styles.css          # Global styles with theming
├── package.json            # Dependencies and scripts
└── README.md              # This file
```

## Development

### Prerequisites

- Node.js 16+ 
- npm or yarn
- Access to Telegram bot for testing

### Installation

```bash
npm install
```

### Development Server

```bash
npm start
```

Runs on http://localhost:3000 (for development outside Telegram)

### Building for Production

```bash
npm run build
```

Creates optimized build in `build/` directory ready for deployment.

### Testing

```bash
npm test
```

## Deployment

### Yandex Cloud Object Storage

1. Build the production bundle:
   ```bash
   npm run build
   ```

2. Upload the `build/` directory contents to Yandex Cloud Object Storage bucket

3. Configure bucket for static website hosting

4. Set HTTPS domain for the bucket

5. Update bot handler with the new URL

### Environment Variables

The app automatically detects the environment:
- **Production**: Must be opened from Telegram
- **Development**: Can run standalone with mock data

## Integration with Backend API

The app communicates with the Go backend API deployed as Yandex Cloud Functions:

- **Endpoint**: `/api/user/tags`
- **Authentication**: Telegram WebApp initData as Bearer token  
- **Response**: JSON array of user tags with message counts

## Telegram WebApp Features Used

- ✅ Theme integration (colors, dark/light mode)
- ✅ User authentication (initData validation)  
- ✅ Haptic feedback on interactions
- ✅ Main button integration (future)
- ✅ Native mobile feel (viewport, no zoom)
- ✅ Closing confirmation

## Design Guidelines

Follows Telegram Mini-App design principles:
- Native mobile feel with proper viewport settings
- Consistent with Telegram UI patterns
- Responsive design optimized for mobile screens
- Proper loading states and error handling
- Accessible color contrast and typography

## Browser Support

- Modern mobile browsers (iOS Safari, Android Chrome)
- Telegram in-app browser
- Desktop browsers for development

## Future Enhancements

- [ ] Tag editing and deletion
- [ ] Message browsing by tag  
- [ ] Search functionality
- [ ] Export features
- [ ] Offline support with caching
- [ ] Push notifications integration

## Troubleshooting

### Common Issues

1. **404 Error**: Ensure API backend is deployed and URL is correct
2. **Authentication Failed**: Check Telegram WebApp SDK integration
3. **Theme Issues**: Verify CSS variables are being set from Telegram
4. **Mobile Issues**: Check viewport meta tags and touch handling

### Development Tips

- Use browser dev tools mobile emulation
- Test with different Telegram themes (light/dark)
- Verify API endpoints with network inspector
- Use React DevTools for component debugging

## License

Part of the Telegram Content Organizer project.