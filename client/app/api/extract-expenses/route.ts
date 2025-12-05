import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    
    const apiKey = process.env.AI_EXPENSE_API_KEY;
    const apiUrl = 'http://localhost:8000/chat/expense_tracker/extract_expenses';
    
    if (!apiKey) {
      return NextResponse.json(
        { success: false, error: 'API key not configured' },
        { status: 500 }
      );
    }

    const response = await fetch(apiUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': apiKey
      },
      body: JSON.stringify(body)
    });

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Extraction API error:', error);
    return NextResponse.json(
      { success: false, error: 'Failed to extract expenses' },
      { status: 500 }
    );
  }
}
