import React, { useState } from "react";
import { Form, Button, Table, Modal } from 'react-bootstrap';
import { Link } from "react-router-dom";

export default function Apis () {
  const [ show, setShow ] = useState(false);
  function handleClose () {
    setShow(false);
  }
  function doSubmit(e) {
    e.preventDefault();
    e.stopPropagation();
    return false;
  }
  return <section className="apis">
    <div className="api-control">
      <Button variant="primary" onClick={() => {
        setShow(true);
      }}>新建</Button>
    </div>
    <div className="apis-list">
      <Table striped bordered hover >
        <thead>
          <tr>
            <th>#</th>
            <th>名称</th>
            <th>类型</th>
            <th>请求路径</th>
            <th>模块名</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>1</td>
            <td>Mark</td>
            <td>Otto</td>
            <td>@mdo</td>
            <td>@mdo</td>
            <td><Link to="/edit">编辑</Link><div>删除</div></td>
          </tr>
          <tr>
            <td>2</td>
            <td>Jacob</td>
            <td>Thornton</td>
            <td>@fat</td>
            <td>@fat</td>
            <td><div>编辑</div><div>删除</div></td>
          </tr>
          <tr>
            <td>3</td>
            <td>Larry the Bird</td>
            <td>@twitter</td>
            <td>@twitter</td>
            <td>@twitter</td>
            <td><div>编辑</div><div>删除</div></td>
          </tr>
        </tbody>
      </Table>
    </div>
    <Modal
      show={show}
      onHide={handleClose}
      backdrop="static"
      keyboard={false}
    >
      <Modal.Header closeButton>
        <Modal.Title>新建</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form onSubmit={doSubmit}>
          <Form.Group className="mb-3">
            <Form.Control placeholder="模块名称(英文)" required/>
            <Form.Control.Feedback type="invalid">
              请输入模块名称(英文)
            </Form.Control.Feedback>
          </Form.Group>
          <Form.Group className="mb-3">
          <Form.Select aria-label="模块类型" required>
            <option value="api">接口</option>
            <option value="model">模块</option>
          </Form.Select>
            <Form.Control.Feedback type="invalid">
              请选择模块类型
            </Form.Control.Feedback>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Control placeholder="模块路径" required/>
            <Form.Control.Feedback type="invalid">
              请输入模块路径(英文)
            </Form.Control.Feedback>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Control placeholder="模块加载名"  disabled/>
          </Form.Group>
        </Form>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleClose}>
          关闭
        </Button>
        <Button variant="primary">新建</Button>
      </Modal.Footer>
    </Modal>
  </section>
}