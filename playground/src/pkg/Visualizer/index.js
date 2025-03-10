import React, {useEffect, useState} from 'react';
import Graph from "react-graph-vis";
import GraphOptions from "./config";

function Visualizer(props) {

    // data
    const [graph, setGraph] = useState({nodes: [], edges: []});

    function modifyNodes() {
        return new Promise((resolve) => {
            let nodes = []
            for (const index in props.graph.nodes) {
                if (props.graph.nodes[index].type === "operation") {
                    nodes.push({
                        id: props.graph.nodes[index].id,
                        label: props.graph.nodes[index].label === "OPERATION_UNION" ? "or" : "and",
                        group: props.graph.nodes[index].type
                    })
                } else {
                    nodes.push({
                        id: props.graph.nodes[index].id,
                        label: props.graph.nodes[index].label,
                        group: props.graph.nodes[index].type
                    })
                }
            }
            resolve(nodes);
        });
    }

    function modifyEdges() {
        return new Promise((resolve) => {
            let edges = []
            for (const index in props.graph.edges) {
                switch (props.graph.edges[index].from.type) {
                    case "entity":
                        edges.push({
                            from: props.graph.edges[index].from.id,
                            to: props.graph.edges[index].to.id,
                            color: {color: 'rgba(99,24,255,0.4)', inherit: false},
                            dashes: false,
                            arrows: {
                                to: {
                                    enabled: false
                                }
                            }
                        })
                        break
                    case "relation":
                        edges.push({
                            from: props.graph.edges[index].from.id,
                            to: props.graph.edges[index].to.id,
                            color: {color: 'rgba(147,241,238,0.4)', inherit: false},
                            dashes: false
                        })
                        break
                    case "permission":
                        edges.push({
                            from: props.graph.edges[index].from.id,
                            to: props.graph.edges[index].to.id,
                            color: {color: 'rgba(91,204,99,0.4)', inherit: false},
                            dashes: false
                        })
                        break
                    case "operation":
                        if (props.graph.edges[index].from.label === "OPERATION_UNION") {
                            edges.push({
                                from: props.graph.edges[index].from.id,
                                to: props.graph.edges[index].to.id,
                                label: props.graph.edges[index].extra ? "not" : "",
                                color: {color: 'rgba(229,52,114,0.4)', inherit: false},
                                dashes: true,
                            })
                        } else {
                            edges.push({
                                from: props.graph.edges[index].from.id,
                                to: props.graph.edges[index].to.id,
                                label: props.graph.edges[index].extra ? "not" : "",
                                color: {color: 'rgba(229,52,114,0.4)', inherit: false},
                                dashes: false,
                            })
                        }
                        break
                    default:
                        edges.push({
                            from: props.graph.edges[index].from.id,
                            to: props.graph.edges[index].to.id,
                        })
                }
            }
            resolve(edges);
        });
    }

    function makeGraph() {
        modifyNodes().then(n => {
            return n
        }).then(nn => {
            modifyEdges().then(e => {
                return {nodes: nn, edges: e}
            }).then(r => {
                if (r.nodes.length > 0) {
                    setGraph(r)
                }
            })
        })
    }

    useEffect(() => {
        setGraph({nodes: [], edges: []})
        if (props.graph) {
            makeGraph()
        }
    }, [props.graph]);

    const events = {};

    return (
        <div style={{height: "85vh"}}>
            {(graph.nodes.length > 0) &&
                <Graph
                    graph={graph}
                    options={GraphOptions()}
                    events={events}
                />
            }
        </div>
    );
}

export default Visualizer;
