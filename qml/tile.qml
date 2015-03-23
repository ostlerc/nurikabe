import QtQuick 2.0

Rectangle {
    id: tile
    property bool open: true
    property int count: 0
    property int index: 0
    width: 25
    height: 25
    border.width: 5
    border.color: "black"
    color: open ? "white" : "black"
    Text {
        anchors.centerIn: parent
        font.pixelSize: 10
        color: "black"
        visible: count > 0
        text: count
    }
    MouseArea {
        id: mouseArea
        anchors.fill: parent
        hoverEnabled: true
        acceptedButtons: Qt.LeftButton | Qt.RightButton
        onClicked: {
            if (count > 0) return
            open = !open
            window.tileChecked(index);
        }
    }
}
