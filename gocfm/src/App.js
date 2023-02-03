import { Routes, Route } from 'react-router';
import './App.css';

import Login from './page/login';
import List from './page/list';
import Edit from './page/edit';
import Apis from './page/apis';
import Layout from './page/layout';

import { UserContext } from './store/user';
import { useContext } from 'react';

function App() {
  const { state } = useContext(UserContext);

  return (
    <Routes>
      {!state.isLogin ?
        <Route path="/" element={<Login />} />
        :
        <Route path="/" element={<Layout />} >
          <Route index element={<List />} />
          <Route path="/apis" element={<Apis/>}/>
          <Route path="/edit" element={<Edit />} />
        </Route>
      }
    </Routes>
  );
}

export default App;
