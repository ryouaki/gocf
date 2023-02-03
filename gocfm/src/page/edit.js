import React, { useEffect, useRef } from "react";

export default function Edit () {
  let editor = useRef(null)

  useEffect(() => {
    window.monaco.editor.create(editor.current, {
      language: 'javascript'
    });
  })
  return <section className="editor">
    <div>

    </div>
    <div ref={editor} id="editor">
      
    </div>
  </section>
}
