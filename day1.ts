import * as fs from 'fs';
const words = fs.readFileSync('day1input.txt', 'utf-8')

const list = words.split("\n")

const leftPart = []
const rightPart = []

for (let index = 0; index < list.length - 1; index++) {
    const element = list[index];
    const numberTable = element.split("   ")
    leftPart.push(numberTable[0])
    rightPart.push(numberTable[1])
}

// Sort both table

leftPart.sort()
rightPart.sort()

console.log(leftPart)
