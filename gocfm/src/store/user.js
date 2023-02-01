import React, { useReducer } from 'react';

let initState = {msg: "h"};

export const UserContext = React.createContext(initState); 

export function UserProvider({children}) {
  const [state, dispatch] = useReducer(function (state, action) {
    switch (action.type) {
      case 'change':
        state.msg = 'hello'
        return {...state};
      default:
        return state;
    }
  }, initState);
  return <UserContext.Provider value={{state, dispatch}}>
    {children}
  </UserContext.Provider>
}