import React, { useContext } from "react";
import { UserContext } from "../store/user";

export default function Login() {
  const {state, dispatch} = useContext(UserContext);
  return  <div>
      111
      {state.msg}
      <button onClick={() => {
        dispatch({ type: 'change' })
      }}>点</button>
    </div>
}