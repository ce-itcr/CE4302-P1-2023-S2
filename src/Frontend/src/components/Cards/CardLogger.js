import React from "react";

export default function CardLogger({ protocol, logs
}) {
    return (
        <>
            <div className="relative flex flex-col min-w-0 break-words bg-white rounded mb-6 xl:mb-0 shadow-lg">
                <div className="flex-auto p-4">
                    <div className="flex flex-wrap">
                        <div className="relative w-full pr-4 max-w-full flex-grow flex-1">
                            <h5 className="text-blueGray-400 uppercase font-bold text-xs">
                                Protocolo actual: <b style={{color: '#FF6100'}}>{protocol}</b>
                            </h5>
                            <span className="font-semibold text-xl text-blueGray-700">
                                Logger de Transacciones del Bus
                            </span>
                        </div>
                    </div>
                    <div className="flex flex-wrap" style={{ paddingTop: 10 }}>
                        <div style={{ height: 395, width: 545 }}>
                            <div style={{height: 385, width: 545, backgroundColor: '#282c34', color: '#fff', padding: 20, borderRadius: 5, fontFamily: 'monospace', overflow: 'auto'}}>
                                {logs.map((log, index) => (
                                    <pre key={index} style={{margin: 0, whiteSpace: 'pre-wrap', wordWrap: 'break-word', lineHeight: 1.5}}>
                                        <code>{log.Log}</code>
                                    </pre>
                                ))}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
}
