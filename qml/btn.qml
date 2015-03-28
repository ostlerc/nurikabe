import QtQuick 2.0
import QtQuick.Controls 1.0
import QtQuick.Controls.Styles 1.0

Button {
    onClicked: window.level(text)
    property string color: "grey"
    property string file
    property bool completed: false
    height: 45

    style: ButtonStyle {
        label: Text {
          renderType: Text.NativeRendering
          font.pointSize: 21
          color: "black"
          text: control.text
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
                    anchors.bottom: parent.bottom
                    anchors.horizontalCenter: parent.horizontalCenter
                    width: 20
                    height: 20
                    source: "images/star.png"
                    visible: control.completed
                }
            }
        }
    }

}
