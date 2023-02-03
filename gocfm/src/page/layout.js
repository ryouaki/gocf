import React from "react";
import { Outlet } from "react-router-dom";
import { Nav } from 'react-bootstrap';

export default function Layout () {
  return <main className="layout df">
    <header>
      <div className="logo">Go Cloud Function</div>
    </header>
    <section className="body df">
      <div className="menu">
        <Nav className="flex-column">
          <Nav.Link href="/">首页</Nav.Link>
          <Nav.Link href="/apis">API列表</Nav.Link>
        </Nav>
      </div>
      <div>
        <Outlet/>
      </div>
    </section>
  </main>
}