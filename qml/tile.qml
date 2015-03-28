import QtQuick 2.0

Rectangle {
    id: tile
    state: "open"
    width: 40
    height: 40
    border.width: 5
    border.color: "black"
    color: "white"

    property int count: 0
    property int index: 0

    states: [
        State {
            name: "open"
            PropertyChanges { target: tile; color: "white" }
        },
        State {
            name: "closed"
            PropertyChanges { target: tile; color: "black" }
        }
    ]

    transitions: Transition {
        ColorAnimation {  properties: "color"; duration: 150 }
    }

    Text {
        anchors.centerIn: parent
        font.pixelSize: 12
        color: "black"
        visible: count > 0
        text: count
    }

    MouseArea {
        id: mouseArea
        anchors.fill: parent
        acceptedButtons: Qt.LeftButton | Qt.RightButton
        onClicked: {
            if (count > 0) return
            tile.state = tile.state == "open" ? "closed" : "open"
            window.tileChecked(index);
        }
    }
}
