import React, { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router";

export default function Edit () {
  let editor = useRef(null);
  let nav = useNavigate()
  let [edit, setEdit] = useState(null)

  useEffect(() => {
    const edit = window.monaco.editor.create(editor.current, {
      language: 'javascript'
    });
    setEdit(edit);
  }, [])

  function pushCode() {
    console.log(edit.getValue())
  }
  return <section className="editor pr">
    <div className="df">
      <div className="file_name">main.js <span onClick={() => {
        nav(-1);
      }}>X</span></div>
    </div>
    <div ref={editor} id="editor">
    </div>
    <div className="push-btn pa" onClick={() => {
      pushCode()
    }}>发布</div>
  </section>
}
