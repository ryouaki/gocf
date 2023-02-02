import React, { useContext, useState } from "react";
import { Form, Button } from 'react-bootstrap';
import { UserContext } from "../store/user";

export default function Login() {
  const [validated, setValidated] = useState(false);
  const {dispatch} = useContext(UserContext);

  function doSubmit(e) {
    e.preventDefault();
    e.stopPropagation();

    setValidated(false);
    dispatch({type: 'doLogin'})
    return false;
  }
  return <main className="m-login">
    <Form className="login-form" onSubmit={doSubmit} noValidate validated={validated}>
      <div className="login-title">Go Cloud Function</div>
      <Form.Group className="mb-3">
        <Form.Control placeholder="账号(admin)" required/>
        <Form.Control.Feedback type="invalid">
          请输入admin
        </Form.Control.Feedback>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Control type="password" placeholder="密码(123456)"required />
        <Form.Control.Feedback type="invalid">
          请输入123456
        </Form.Control.Feedback>
      </Form.Group>
      <Button variant="primary" type="submit" className="login-submit">
        Submit
      </Button>
    </Form>
  </main>
}