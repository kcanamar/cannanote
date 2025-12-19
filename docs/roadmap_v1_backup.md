# CannaNote Roadmap

## Project Vision
CannaNote aims to be a comprehensive platform for cannabis users to track, share, and discover cannabis experiences. The goal is to create a social community where users can document their unique reactions to different strains and consumption methods.

## Current Status (v1.0)

### ✅ Completed Features
- User authentication (signup/login/logout)
- Basic entry CRUD operations (create, read, update, delete)
- Entry model with strain, type, amount, consumption method, and description
- Session-based authentication with MongoDB storage
- Basic like/favorite functionality
- EJS templating with Tailwind CSS styling
- Responsive design foundation
- Basic HTMX integration for delete confirmation

### 🔧 Technical Debt & Bug Fixes
- **Like Button Limitation**: Currently allows unlimited likes per user (needs single-use constraint)
- **Onboarding Documentation**: Missing user onboarding flow documentation
- **Error Handling**: Improve error handling across all controllers
- **Validation**: Add form validation for entry creation and editing

## Roadmap

### Phase 1: Core UX Improvements (Next 2-4 weeks)
**Priority: High**

#### Enhanced HTMX Integration
- [ ] Implement "click to edit" functionality in entry views
- [ ] Move edit functionality to show page with inline editing
- [ ] Add real-time like updates without page refresh
- [ ] Implement dynamic entry filtering/sorting

#### User Experience Polish
- [ ] Add form validation with user-friendly error messages
- [ ] Implement loading states for all async operations
- [ ] Add confirmation modals for destructive actions
- [ ] Improve responsive design for mobile devices

#### Data Quality
- [ ] Standardize strain names with autocomplete suggestions
- [ ] Add consumption method validation and suggestions
- [ ] Implement entry tagging system
- [ ] Add entry photos/images capability

### Phase 2: Social Features (1-2 months)
**Priority: Medium-High**

#### User Profiles & Social Interaction
- [ ] User profile pages with entry history
- [ ] Follow/unfollow user functionality
- [ ] Comment system on entries
- [ ] Entry sharing capabilities
- [ ] User activity feeds

#### Discovery & Search
- [ ] Advanced search and filtering by strain, type, effects
- [ ] Trending strains and popular entries
- [ ] Recommendation system based on user preferences
- [ ] Geographic strain availability tracking

### Phase 3: Advanced Features (2-3 months)
**Priority: Medium**

#### Age Verification System
- [ ] Age verification gate for new users
- [ ] Compliance with local cannabis laws
- [ ] Terms of service and privacy policy

#### Medical Cannabis Focus
- [ ] Medical patient verification system
- [ ] Symptom tracking integration
- [ ] Dosage recommendations based on medical needs
- [ ] Healthcare provider integration potential

#### Analytics & Insights
- [ ] Personal consumption analytics dashboard
- [ ] Strain effectiveness tracking over time
- [ ] Export personal data functionality
- [ ] Integration with health tracking apps

### Phase 4: Platform Expansion (3-6 months)
**Priority: Low-Medium**

#### Notifications System
- [ ] Email notifications for follows/comments
- [ ] In-app notification system
- [ ] Push notifications for mobile web

#### API Development
- [ ] RESTful API for third-party integrations
- [ ] Mobile app development preparation
- [ ] Dispensary integration possibilities

#### Community Features
- [ ] User-generated strain reviews and ratings
- [ ] Community challenges and events
- [ ] Educational content integration
- [ ] Expert contributor program

## Technical Considerations

### Performance Optimizations
- [ ] Implement database indexing for search performance
- [ ] Add caching layer for frequently accessed data
- [ ] Optimize image storage and delivery
- [ ] Implement pagination for large datasets

### Security & Compliance
- [ ] Implement rate limiting
- [ ] Add CSRF protection
- [ ] Security audit and penetration testing
- [ ] GDPR compliance for international users

### Infrastructure
- [ ] Set up proper staging environment
- [ ] Implement automated testing suite
- [ ] Add monitoring and logging
- [ ] Database backup and recovery procedures

## Success Metrics

### User Engagement
- Monthly active users
- Entry creation frequency
- User retention rates
- Social interaction metrics (likes, follows, comments)

### Platform Health
- Page load times
- Error rates
- User satisfaction scores
- Mobile usage statistics

## Contributing
When working on roadmap items:
1. Create feature branches from `main`
2. Follow existing code conventions
3. Add appropriate tests for new functionality
4. Update documentation as needed
5. Create pull requests for review

## Notes
- Priority levels may shift based on user feedback and usage analytics
- Timeline estimates are subject to change based on resource availability
- Community feedback will heavily influence feature prioritization