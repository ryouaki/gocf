import { Routes, Route } from 'react-router';
import './App.css';

import Login from './page/login';
import List from './page/list';
import Edit from './page/edit';



function App() {
  return (
      <Routes>
          <Route path="/" element={<Login />} />
          <Route path="/list" element={<List />} />
          <Route path="/edit" element={<Edit />} />
      </Routes>
    
  );
}

export default App;
