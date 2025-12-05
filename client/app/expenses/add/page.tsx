'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuthStore } from '@/store/authStore';
import api from '@/lib/api';
import toast from 'react-hot-toast';
import { Category } from '@/types';

interface ExtractedExpense {
  amount: number;
  category: string;
  category_id: string;
  description: string;
  date: string | null;
  merchant: string | null;
}

interface ParsedExpense {
  name: string;
  category_id: number;
  unit: number;
  per_unit_cost: number;
}

export default function AddExpensePage() {
  const router = useRouter();
  const { isAuthenticated, loadFromStorage } = useAuthStore();
  const [categories, setCategories] = useState<Category[]>([]);
  const [name, setName] = useState('');
  const [categoryId, setCategoryId] = useState('');
  const [unit, setUnit] = useState('');
  const [perUnitCost, setPerUnitCost] = useState('');
  const [total, setTotal] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [showNewCategory, setShowNewCategory] = useState(false);
  const [newCategoryName, setNewCategoryName] = useState('');
  
  // AI Extraction states
  const [showAiModal, setShowAiModal] = useState(false);
  const [expenseText, setExpenseText] = useState('');
  const [isExtracting, setIsExtracting] = useState(false);
  const [extractedExpenses, setExtractedExpenses] = useState<ParsedExpense[]>([]);
  const [showExtractedList, setShowExtractedList] = useState(false);

  useEffect(() => {
    loadFromStorage();
  }, [loadFromStorage]);

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/auth/login');
      return;
    }
    fetchCategories();
  }, [isAuthenticated, router]);

  useEffect(() => {
    const unitNum = parseFloat(unit) || 0;
    const costNum = parseFloat(perUnitCost) || 0;
    setTotal(unitNum * costNum);
  }, [unit, perUnitCost]);

  const fetchCategories = async () => {
    try {
      const response = await api.get('/categories');
      setCategories(response.data.data);
    } catch (error) {
      toast.error('Failed to load categories');
    }
  };

  const handleCreateCategory = async () => {
    if (!newCategoryName.trim()) {
      toast.error('Category name is required');
      return;
    }

    try {
      const response = await api.post('/categories', {
        name: newCategoryName,
      });
      const newCategory = response.data.data;
      setCategories([...categories, newCategory]);
      setCategoryId(newCategory.id.toString());
      setNewCategoryName('');
      setShowNewCategory(false);
      toast.success('Category created successfully');
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to create category');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!categoryId || !unit || !perUnitCost) {
      toast.error('Please fill in all fields');
      return;
    }

    setIsLoading(true);
    try {
      await api.post('/expenses', {
        name,
        category_id: parseInt(categoryId),
        unit: parseFloat(unit),
        per_unit_cost: parseFloat(perUnitCost),
        expense_date: new Date().toISOString(),
      });
      toast.success('Expense added successfully!');
      router.push('/dashboard');
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to add expense');
    } finally {
      setIsLoading(false);
    }
  };

  const mapCategoryNameToId = (categoryName: string): number => {
    const mapping: Record<string, string> = {
      'groceries': 'Food & Dining',
      'dining': 'Food & Dining',
      'food': 'Food & Dining',
      'transportation': 'Transportation',
      'shopping': 'Shopping',
      'entertainment': 'Entertainment',
      'bills': 'Bills & Utilities',
      'utilities': 'Bills & Utilities',
      'healthcare': 'Healthcare',
      'education': 'Education',
      'personal care': 'Personal Care',
      'travel': 'Travel',
    };

    const normalizedName = categoryName.toLowerCase();
    const mappedName = mapping[normalizedName] || categoryName;
    
    const category = categories.find(
      c => c.name.toLowerCase() === mappedName.toLowerCase()
    );
    
    return category?.id || categories.find(c => c.name === 'Other')?.id || 10;
  };

  const handleExtractExpenses = async () => {
    if (!expenseText.trim()) {
      toast.error('Please enter expense text');
      return;
    }

    setIsExtracting(true);
    try {
      const categoryList = categories.map(cat => ({
        category_id: `cat_${cat.id.toString().padStart(3, '0')}`,
        name: cat.name.toLowerCase(),
        is_default: cat.is_default
      }));

      const response = await api.post('/expenses/extract', {
        input_data: {
          paragraph: expenseText,
          categories: categoryList
        },
        conversation_history: []
      });

      const data = response.data;
      
      if (data.success && data.output_data?.expenses) {
        const parsed: ParsedExpense[] = data.output_data.expenses.map((exp: ExtractedExpense) => ({
          name: exp.description || 'Expense',
          category_id: mapCategoryNameToId(exp.category),
          unit: 1,
          per_unit_cost: exp.amount
        }));

        setExtractedExpenses(parsed);
        setShowExtractedList(true);
        setShowAiModal(false);
        toast.success(`Extracted ${parsed.length} expenses!`);
      } else {
        toast.error('Failed to extract expenses');
      }
    } catch (error) {
      console.error('Extraction error:', error);
      toast.error('Failed to extract expenses');
    } finally {
      setIsExtracting(false);
    }
  };

  const handleAddAllExtracted = async () => {
    if (extractedExpenses.length === 0) {
      toast.error('No expenses to add');
      return;
    }

    setIsLoading(true);
    let successCount = 0;
    let errorCount = 0;

    for (const expense of extractedExpenses) {
      try {
        await api.post('/expenses', {
          ...expense,
          expense_date: new Date().toISOString(),
        });
        successCount++;
      } catch (error) {
        errorCount++;
      }
    }

    setIsLoading(false);
    
    if (successCount > 0) {
      toast.success(`Added ${successCount} expense(s) successfully!`);
      if (errorCount === 0) {
        router.push('/dashboard');
      }
    }
    
    if (errorCount > 0) {
      toast.error(`Failed to add ${errorCount} expense(s)`);
    }
  };

  const updateExtractedExpense = (index: number, field: keyof ParsedExpense, value: any) => {
    const updated = [...extractedExpenses];
    updated[index] = { ...updated[index], [field]: value };
    setExtractedExpenses(updated);
  };

  const removeExtractedExpense = (index: number) => {
    setExtractedExpenses(extractedExpenses.filter((_, i) => i !== index));
  };

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <Link href="/dashboard" className="text-xl font-bold text-gray-900">
                MinFlow
              </Link>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="bg-white shadow rounded-lg">
            <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg leading-6 font-medium text-gray-900">
                    Add New Expense
                  </h3>
                  <p className="mt-1 text-sm text-gray-500">
                    Record your expense details below
                  </p>
                </div>
                <button
                  type="button"
                  onClick={() => setShowAiModal(true)}
                  className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-purple-600 hover:bg-purple-700"
                >
                  <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                  AI Extract
                </button>
              </div>
            </div>

            {!showExtractedList && (
              <form onSubmit={handleSubmit} className="px-4 py-5 sm:p-6 space-y-6">
              {/* Expense Name */}
              <div>
                <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                  Expense Name
                </label>
                <input
                  type="text"
                  id="name"
                  required
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  placeholder="e.g., Groceries, Fuel, etc."
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </div>

              {/* Category Selection */}
              <div>
                <label htmlFor="category" className="block text-sm font-medium text-gray-700">
                  Category
                </label>
                <div className="mt-1 flex space-x-2">
                  <select
                    id="category"
                    required
                    className="block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                    value={categoryId}
                    onChange={(e) => setCategoryId(e.target.value)}
                  >
                    <option value="">Select a category</option>
                    {categories.map((category) => (
                      <option key={category.id} value={category.id}>
                        {category.name}
                      </option>
                    ))}
                  </select>
                  <button
                    type="button"
                    onClick={() => setShowNewCategory(!showNewCategory)}
                    className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
                  >
                    New
                  </button>
                </div>

                {/* New Category Form */}
                {showNewCategory && (
                  <div className="mt-2 flex space-x-2">
                    <input
                      type="text"
                      className="block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                      placeholder="New category name"
                      value={newCategoryName}
                      onChange={(e) => setNewCategoryName(e.target.value)}
                    />
                    <button
                      type="button"
                      onClick={handleCreateCategory}
                      className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
                    >
                      Create
                    </button>
                    <button
                      type="button"
                      onClick={() => {
                        setShowNewCategory(false);
                        setNewCategoryName('');
                      }}
                      className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
                    >
                      Cancel
                    </button>
                  </div>
                )}
              </div>

              {/* Unit */}
              <div>
                <label htmlFor="unit" className="block text-sm font-medium text-gray-700">
                  Quantity/Unit
                </label>
                <input
                  type="number"
                  id="unit"
                  required
                  step="0.01"
                  min="0.01"
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  placeholder="e.g., 2.5"
                  value={unit}
                  onChange={(e) => setUnit(e.target.value)}
                />
              </div>

              {/* Per Unit Cost */}
              <div>
                <label htmlFor="cost" className="block text-sm font-medium text-gray-700">
                  Per Unit Cost
                </label>
                <input
                  type="number"
                  id="cost"
                  required
                  step="0.01"
                  min="0.01"
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  placeholder="e.g., 100.00"
                  value={perUnitCost}
                  onChange={(e) => setPerUnitCost(e.target.value)}
                />
              </div>

              {/* Total (Read-only) */}
              <div>
                <label htmlFor="total" className="block text-sm font-medium text-gray-700">
                  Total Amount
                </label>
                <input
                  type="text"
                  id="total"
                  readOnly
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 bg-gray-50 text-gray-900 sm:text-sm font-semibold"
                  value={`₹ ${total.toFixed(2)}`}
                />
              </div>

              {/* Submit Button */}
              <div className="flex justify-end space-x-3">
                <Link
                  href="/dashboard"
                  className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
                >
                  Cancel
                </Link>
                <button
                  type="submit"
                  disabled={isLoading}
                  className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                >
                  {isLoading ? 'Adding...' : 'Add Expense'}
                </button>
              </div>
            </form>
            )}

            {/* Extracted Expenses List */}
            {showExtractedList && extractedExpenses.length > 0 && (
              <div className="px-4 py-5 sm:p-6">
                <div className="flex justify-between items-center mb-4">
                  <h4 className="text-lg font-medium text-gray-900">
                    Review Extracted Expenses ({extractedExpenses.length})
                  </h4>
                  <button
                    type="button"
                    onClick={() => setShowExtractedList(false)}
                    className="text-sm text-gray-600 hover:text-gray-900"
                  >
                    ← Back to manual entry
                  </button>
                </div>

                <div className="space-y-4">
                  {extractedExpenses.map((expense, index) => (
                    <div key={index} className="border border-gray-200 rounded-lg p-4">
                      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                        <div>
                          <label className="block text-xs font-medium text-gray-700 mb-1">
                            Name
                          </label>
                          <input
                            type="text"
                            className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm text-gray-900"
                            value={expense.name}
                            onChange={(e) => updateExtractedExpense(index, 'name', e.target.value)}
                          />
                        </div>
                        <div>
                          <label className="block text-xs font-medium text-gray-700 mb-1">
                            Category
                          </label>
                          <select
                            className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm text-gray-900"
                            value={expense.category_id}
                            onChange={(e) => updateExtractedExpense(index, 'category_id', parseInt(e.target.value))}
                          >
                            {categories.map((cat) => (
                              <option key={cat.id} value={cat.id}>
                                {cat.name}
                              </option>
                            ))}
                          </select>
                        </div>
                        <div>
                          <label className="block text-xs font-medium text-gray-700 mb-1">
                            Unit
                          </label>
                          <input
                            type="number"
                            step="0.01"
                            min="0.01"
                            className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm text-gray-900"
                            value={expense.unit}
                            onChange={(e) => updateExtractedExpense(index, 'unit', parseFloat(e.target.value))}
                          />
                        </div>
                        <div>
                          <label className="block text-xs font-medium text-gray-700 mb-1">
                            Cost
                          </label>
                          <div className="flex items-center space-x-2">
                            <input
                              type="number"
                              step="0.01"
                              min="0.01"
                              className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm text-gray-900"
                              value={expense.per_unit_cost}
                              onChange={(e) => updateExtractedExpense(index, 'per_unit_cost', parseFloat(e.target.value))}
                            />
                            <button
                              type="button"
                              onClick={() => removeExtractedExpense(index)}
                              className="text-red-600 hover:text-red-800"
                            >
                              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </div>
                      <div className="mt-2 text-sm text-gray-600">
                        Total: ₹{(expense.unit * expense.per_unit_cost).toFixed(2)}
                      </div>
                    </div>
                  ))}
                </div>

                <div className="mt-6 flex justify-end space-x-3">
                  <button
                    type="button"
                    onClick={() => {
                      setExtractedExpenses([]);
                      setShowExtractedList(false);
                    }}
                    className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
                  >
                    Cancel
                  </button>
                  <button
                    type="button"
                    onClick={handleAddAllExtracted}
                    disabled={isLoading || extractedExpenses.length === 0}
                    className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50"
                  >
                    {isLoading ? 'Adding...' : `Add All ${extractedExpenses.length} Expense(s)`}
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </main>

      {/* AI Extraction Modal */}
      {showAiModal && (
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-2xl w-full p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">
              AI Expense Extraction
            </h3>
            <p className="text-sm text-gray-600 mb-4">
              Paste or type your expenses in natural language. For example: &quot;tin hali kola 30 taka, 600g murgi 200 taka, riksha vara 50&quot;
            </p>
            <textarea
              className="w-full h-40 border border-gray-300 rounded-md p-3 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              placeholder="Enter your expenses here..."
              value={expenseText}
              onChange={(e) => setExpenseText(e.target.value)}
            />
            <div className="mt-4 flex justify-end space-x-3">
              <button
                type="button"
                onClick={() => {
                  setShowAiModal(false);
                  setExpenseText('');
                }}
                className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={handleExtractExpenses}
                disabled={isExtracting || !expenseText.trim()}
                className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-purple-600 hover:bg-purple-700 disabled:opacity-50"
              >
                {isExtracting ? 'Extracting...' : 'Extract Expenses'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
