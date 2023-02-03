import React, { useReducer } from 'react';

const initState = { isLogin: true };

export const UserContext = React.createContext(initState);

export function UserProvider({ children }) {
  const [state, dispatch] = useReducer(function (state, action) {
    switch (action.type) {
      case 'doLogin':
        state.isLogin = true;
        return { ...state };
      case 'doLogout':
        state.isLogin = false;
        return { ...state };
      default:
        return state;
    }
  }, initState);
  return <UserContext.Provider value={{ state, dispatch }}>
    {children}
  </UserContext.Provider>
}