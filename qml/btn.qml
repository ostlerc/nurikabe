import QtQuick 2.0
import QtQuick.Controls 1.0
import QtQuick.Controls.Styles 1.0

Button {
    onClicked: window.onBtnClicked(data)
    property string color: "lightsteelblue"
    property string data
    property bool completed: false
    property bool showstar: false
    property bool alignCenter: false

    style: ButtonStyle {
        label: Text {
            renderType: Text.NativeRendering
            font.pointSize: 20
            color: "black"
            text: control.text

            verticalAlignment: Text.AlignVCenter
            horizontalAlignment: control.alignCenter ? Text.AlignHCenter : Text.AlignLeft
            anchors.fill: parent
        }

        background: Component {
            Rectangle {
                border.width: 1
                radius: 5
                gradient: Gradient {
                    GradientStop { position: 0 ; color: control.pressed ? Qt.darker(control.color) : control.color }
                    GradientStop { position: 1 ; color: control.pressed ? Qt.darker(control.color) : control.color }
                }
                anchors.fill: parent

                Image {
                    anchors.right: parent.right
                    anchors.verticalCenter: parent.verticalCenter
                    width: 20
                    height: 20
                    source: control.completed ? "images/star.png" : "images/emptystar.png"
                    visible: showstar
                }
            }
        }
    }
}
