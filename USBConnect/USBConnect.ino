// include the library code:
#include <LiquidCrystal.h>

// temperature measurement string format: machine_name : temperature
String currTemp = "";
String lastTemp = "";


// initialize the library with the numbers of the interface pins
// https://learn.adafruit.com/adafruit-arduino-lesson-11-lcd-displays-1/other-things-to-do?view=all#arduino-code
LiquidCrystal lcd(7, 8, 9, 10, 11, 12);

void setup() {
  // listen for temperature measurements
  Serial.begin(9600);
  while (!Serial) {
    ; // wait for serial port to connect. Needed for native USB port only
  }

  // set up the LCD's number of columns and rows:
  lcd.begin(16, 2);
  
}
  
void loop() {

  currTemp = Serial.readStringUntil('\n');

  // set the cursor to column 0, line 1
  // (note: line 1 is the second row, since counting begins with 0):
  lcd.setCursor(0, 0);
  lcd.print(currTemp);
  
  lcd.setCursor(0, 1);
  lcd.print(lastTemp);

  if(currTemp != lastTemp) {
    lastTemp = currTemp;
  }

  delay(5000);

  lcd.clear();

}

