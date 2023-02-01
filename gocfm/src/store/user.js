import React from 'react';

const UserContext = React.createContext({});

export default function UserProvider({childrens}) {
  return <UserContext.Provider>
    {childrens}
  </UserContext.Provider>
}