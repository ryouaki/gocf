import { Routes, Route } from 'react-router';
import './App.css';

import Login from './page/login';
import List from './page/list';
import Edit from './page/edit';
import Home from './page/home';

import { UserContext } from './store/user';
import { useContext } from 'react';

function App() {
  const { state } = useContext(UserContext);

  return (
    <Routes>
      {!state.isLogin ? <Route path="/" element={<Login />} /> : <Route path="/" element={<Home />} />}
      <Route path="/list" element={<List />} />
      <Route path="/edit" element={<Edit />} />
    </Routes>
  );
}

export default App;
