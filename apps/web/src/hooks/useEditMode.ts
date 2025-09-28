import { useState, useCallback } from 'react';

export type ViewMode = 'view' | 'edit';

export interface UseEditModeReturn {
  mode: ViewMode;
  isEditMode: boolean;
  enterEditMode: () => void;
  exitEditMode: () => void;
  setMode: (mode: ViewMode) => void;
}

export function useEditMode(initialMode: ViewMode = 'view'): UseEditModeReturn {
  const [mode, setModeState] = useState<ViewMode>(initialMode);

  const isEditMode = mode === 'edit';

  const enterEditMode = useCallback(() => {
    setModeState('edit');
  }, []);

  const exitEditMode = useCallback(() => {
    setModeState('view');
  }, []);

  const setMode = useCallback((newMode: ViewMode) => {
    setModeState(newMode);
  }, []);

  return {
    mode,
    isEditMode,
    enterEditMode,
    exitEditMode,
    setMode,
  };
}
