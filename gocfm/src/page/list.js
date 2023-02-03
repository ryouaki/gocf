import React from "react";
import { Card } from 'react-bootstrap';

export default function List() {
  const items = new Array(20).fill(1)
  return <section className="list df">
    {items.map(() => {
      return <Card style={{ width: '18rem' }} className="pod-card">
        <Card.Body>
          <Card.Title>10.198.2.1</Card.Title>
          <Card.Body className="pod-wrap">
            <div>最后更新: 2021.01.01 10:10:10</div>
            <div>CPU: 32%,MEM: 200/1000</div>
          </Card.Body>
          <Card.Body className="pod-wrap df">
            <div className="pod-green"></div><div className="warning">正常</div>
          </Card.Body>
        </Card.Body>
      </Card>
    })}
  </section>
}