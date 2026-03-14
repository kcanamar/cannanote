// CannaNote Copy Service
// Hexagonal architecture: Port + Adapters for voice-based copy
// Supports: holistic (professional), westcoast (recreational), ganjier (connoisseur)

const CopyService = (function() {
  // Default voice
  const DEFAULT_VOICE = 'holistic';

  // Current voice (loaded from preferences)
  let currentVoice = DEFAULT_VOICE;

  // Voice adapters - each implements the same copy interface
  const adapters = {
    // Professional/Medical voice - direct, precise, clinical but human
    holistic: {
      // Dashboard
      'dashboard.logSession.title': 'Log Session',
      'dashboard.logSession.subtitle': 'Record your experience',
      'dashboard.stats.mostUsed': 'Most Used',
      'dashboard.stats.bestMood': 'Best for Mood',
      'dashboard.stats.collection': 'Your Collection',
      'dashboard.recentSessions': 'Recent Sessions',
      'dashboard.noSessions': 'No sessions recorded. Begin logging to track your experiences.',

      // Session logging
      'sessionLog.title': 'New Session',
      'sessionLog.step1.title': 'Product',
      'sessionLog.step1.subtitle': 'Select or add product',
      'sessionLog.step2.title': 'Consumption',
      'sessionLog.step2.subtitle': 'Method and amount',
      'sessionLog.step3.title': 'Pre-Session',
      'sessionLog.step3.subtitle': 'Current state assessment',
      'sessionLog.product.strainLabel': 'Strain Name',
      'sessionLog.product.strainPlaceholder': 'Enter strain name',
      'sessionLog.product.producerLabel': 'Producer',
      'sessionLog.product.producerPlaceholder': 'Enter producer (optional)',
      'sessionLog.method.label': 'Consumption Method',
      'sessionLog.amount.label': 'Amount',
      'sessionLog.feeling.mood': 'Mood',
      'sessionLog.feeling.mind': 'Mind',
      'sessionLog.feeling.body': 'Body',
      'sessionLog.feeling.prompt': 'Rate your current state',
      'sessionLog.submit': 'Begin Session',
      'sessionLog.checkin.enable': 'Enable check-ins',
      'sessionLog.checkin.interval': 'Check-in interval',

      // Check-in
      'checkin.title': 'Check-in',
      'checkin.howFeeling': 'Current state assessment',
      'checkin.intensity': 'Effect Intensity',
      'checkin.expectations': 'Does this match your expectations?',
      'checkin.expectations.yes': 'Yes',
      'checkin.expectations.somewhat': 'Somewhat',
      'checkin.expectations.no': 'No',
      'checkin.notes': 'Notes (Optional)',
      'checkin.notesPlaceholder': 'Document any observations',
      'checkin.submit': 'Save Check-in',
      'checkin.endSession': 'End Session',

      // Sessions list
      'sessions.title': 'Sessions',
      'sessions.active': 'Active',
      'sessions.completed': 'Completed',
      'sessions.empty': 'No sessions recorded.',
      'sessions.startFirst': 'Record your first session to begin tracking.',

      // Settings
      'settings.title': 'Settings',
      'settings.subtitle': 'Configure preferences',
      'settings.voice.title': 'Voice',
      'settings.voice.subtitle': 'Select communication style',

      // General
      'general.save': 'Save',
      'general.cancel': 'Cancel',
      'general.delete': 'Delete',
      'general.back': 'Back',
      'general.next': 'Next',
      'general.loading': 'Loading...',
      'general.error': 'An error occurred'
    },

    // West Coast/California voice - laid-back, recreational, familiar
    westcoast: {
      // Dashboard
      'dashboard.logSession.title': 'Log a Sesh',
      'dashboard.logSession.subtitle': 'What\'s good today?',
      'dashboard.stats.mostUsed': 'Go-To Strain',
      'dashboard.stats.bestMood': 'Mood Booster',
      'dashboard.stats.collection': 'Your Stash',
      'dashboard.recentSessions': 'Recent Seshes',
      'dashboard.noSessions': 'No seshes yet. Light one up and start tracking!',

      // Session logging
      'sessionLog.title': 'New Sesh',
      'sessionLog.step1.title': 'The Goods',
      'sessionLog.step1.subtitle': 'What are you smoking?',
      'sessionLog.step2.title': 'How',
      'sessionLog.step2.subtitle': 'Method and dose',
      'sessionLog.step3.title': 'Vibe Check',
      'sessionLog.step3.subtitle': 'How you feeling right now?',
      'sessionLog.product.strainLabel': 'Strain',
      'sessionLog.product.strainPlaceholder': 'What strain?',
      'sessionLog.product.producerLabel': 'From',
      'sessionLog.product.producerPlaceholder': 'Where\'d you get it?',
      'sessionLog.method.label': 'How you taking it?',
      'sessionLog.amount.label': 'How much?',
      'sessionLog.feeling.mood': 'Mood',
      'sessionLog.feeling.mind': 'Head',
      'sessionLog.feeling.body': 'Body',
      'sessionLog.feeling.prompt': 'How you vibin\'?',
      'sessionLog.submit': 'Let\'s Go',
      'sessionLog.checkin.enable': 'Remind me to check in',
      'sessionLog.checkin.interval': 'Remind me in',

      // Check-in
      'checkin.title': 'Check-in',
      'checkin.howFeeling': 'How you feelin\'?',
      'checkin.intensity': 'How strong?',
      'checkin.expectations': 'Is it hitting like you expected?',
      'checkin.expectations.yes': 'Yeah',
      'checkin.expectations.somewhat': 'Kinda',
      'checkin.expectations.no': 'Nah',
      'checkin.notes': 'Any thoughts?',
      'checkin.notesPlaceholder': 'What\'s on your mind?',
      'checkin.submit': 'Save',
      'checkin.endSession': 'Done for now',

      // Sessions list
      'sessions.title': 'Your Seshes',
      'sessions.active': 'Going',
      'sessions.completed': 'Done',
      'sessions.empty': 'No seshes yet.',
      'sessions.startFirst': 'Start your first sesh!',

      // Settings
      'settings.title': 'Settings',
      'settings.subtitle': 'Make it yours',
      'settings.voice.title': 'Vibe',
      'settings.voice.subtitle': 'How should we talk to you?',

      // General
      'general.save': 'Save',
      'general.cancel': 'Nevermind',
      'general.delete': 'Delete',
      'general.back': 'Back',
      'general.next': 'Next',
      'general.loading': 'Loading...',
      'general.error': 'Something went wrong'
    },

    // Ganjier/Connoisseur voice - craft appreciation, sommelier energy
    ganjier: {
      // Dashboard
      'dashboard.logSession.title': 'Document Experience',
      'dashboard.logSession.subtitle': 'Curate your journey',
      'dashboard.stats.mostUsed': 'Signature Cultivar',
      'dashboard.stats.bestMood': 'Finest for Elevation',
      'dashboard.stats.collection': 'Your Cellar',
      'dashboard.recentSessions': 'Recent Experiences',
      'dashboard.noSessions': 'Your journal awaits. Document your first cultivar experience.',

      // Session logging
      'sessionLog.title': 'New Experience',
      'sessionLog.step1.title': 'The Cultivar',
      'sessionLog.step1.subtitle': 'Select your flower',
      'sessionLog.step2.title': 'Consumption',
      'sessionLog.step2.subtitle': 'Method and portion',
      'sessionLog.step3.title': 'Baseline',
      'sessionLog.step3.subtitle': 'Current disposition',
      'sessionLog.product.strainLabel': 'Cultivar',
      'sessionLog.product.strainPlaceholder': 'Enter cultivar name',
      'sessionLog.product.producerLabel': 'Cultivator',
      'sessionLog.product.producerPlaceholder': 'Farm or producer',
      'sessionLog.method.label': 'Consumption Method',
      'sessionLog.amount.label': 'Portion',
      'sessionLog.feeling.mood': 'Disposition',
      'sessionLog.feeling.mind': 'Mental Clarity',
      'sessionLog.feeling.body': 'Physical State',
      'sessionLog.feeling.prompt': 'Assess your current state',
      'sessionLog.submit': 'Begin Experience',
      'sessionLog.checkin.enable': 'Schedule reflections',
      'sessionLog.checkin.interval': 'Reflection interval',

      // Check-in
      'checkin.title': 'Reflection',
      'checkin.howFeeling': 'Present state assessment',
      'checkin.intensity': 'Effect Magnitude',
      'checkin.expectations': 'Does this express as expected?',
      'checkin.expectations.yes': 'Precisely',
      'checkin.expectations.somewhat': 'Partially',
      'checkin.expectations.no': 'Unexpected',
      'checkin.notes': 'Tasting Notes',
      'checkin.notesPlaceholder': 'Describe the expression, terroir, effects...',
      'checkin.submit': 'Record Reflection',
      'checkin.endSession': 'Conclude Experience',

      // Sessions list
      'sessions.title': 'Experience Journal',
      'sessions.active': 'In Progress',
      'sessions.completed': 'Archived',
      'sessions.empty': 'Your journal is empty.',
      'sessions.startFirst': 'Begin documenting your cultivar experiences.',

      // Settings
      'settings.title': 'Preferences',
      'settings.subtitle': 'Customize your experience',
      'settings.voice.title': 'Voice',
      'settings.voice.subtitle': 'Select your preferred tone',

      // General
      'general.save': 'Preserve',
      'general.cancel': 'Discard',
      'general.delete': 'Remove',
      'general.back': 'Return',
      'general.next': 'Continue',
      'general.loading': 'Preparing...',
      'general.error': 'An issue occurred'
    }
  };

  // Get copy by key
  function get(key, fallback) {
    const adapter = adapters[currentVoice] || adapters[DEFAULT_VOICE];
    return adapter[key] || fallback || key;
  }

  // Get all copy for current voice (for bulk operations)
  function getAll() {
    return { ...adapters[currentVoice] } || { ...adapters[DEFAULT_VOICE] };
  }

  // Set current voice
  function setVoice(voice) {
    if (adapters[voice]) {
      currentVoice = voice;
      console.log('[CopyService] Voice set to:', voice);
      return true;
    }
    console.warn('[CopyService] Unknown voice:', voice);
    return false;
  }

  // Get current voice
  function getVoice() {
    return currentVoice;
  }

  // Get available voices
  function getVoices() {
    return Object.keys(adapters);
  }

  // Initialize from preferences
  async function init() {
    if (typeof SessionDB !== 'undefined') {
      try {
        const voice = await SessionDB.getPreference('voice', DEFAULT_VOICE);
        setVoice(voice);
      } catch (err) {
        console.error('[CopyService] Failed to load voice preference:', err);
      }
    }

    // Listen for voice changes from Settings
    window.addEventListener('voiceChanged', (e) => {
      if (e.detail && e.detail.voice) {
        setVoice(e.detail.voice);
        // Trigger UI update
        window.dispatchEvent(new CustomEvent('copyUpdated'));
      }
    });
  }

  // Apply copy to elements with data-copy attribute
  function applyToDOM() {
    document.querySelectorAll('[data-copy]').forEach(el => {
      const key = el.dataset.copy;
      const text = get(key);
      if (text && text !== key) {
        if (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA') {
          if (el.placeholder) {
            el.placeholder = text;
          }
        } else {
          el.textContent = text;
        }
      }
    });
  }

  // Public API
  return {
    init,
    get,
    getAll,
    setVoice,
    getVoice,
    getVoices,
    applyToDOM
  };
})();

// Initialize on load
if (typeof window !== 'undefined') {
  window.CopyService = CopyService;

  document.addEventListener('DOMContentLoaded', async () => {
    await CopyService.init();
    CopyService.applyToDOM();
  });

  // Re-apply when voice changes
  window.addEventListener('copyUpdated', () => {
    CopyService.applyToDOM();
  });
}
