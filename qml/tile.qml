import QtQuick 2.0

Rectangle {
    id: tile
    property int type: 0
    property int index: 0
    property int count: 0
    width: 25
    height: 25
    color: {
        if (type == 0) // Open
            return "white"
        return "black" // Closed
    }
    border.color: {
        return "black"
    }
    Text {
        anchors.centerIn: parent
        font.pixelSize: 10
        color: "black"
        visible: count != 0
        text: count
    }
    border.width: 5
    MouseArea {
        id: mouseArea
        anchors.fill: parent
        hoverEnabled: true
        acceptedButtons: Qt.LeftButton | Qt.RightButton
        onClicked: {
            if (count != 0) {
                return
            }

            if (type == 0) {
                type = 1
            } else {
                type = 0
            }
            window.tileChecked(index);
        }
    }
}
